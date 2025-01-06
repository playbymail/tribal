// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package pigeon

import (
	"log"
	"time"
)

const (
	// CR and LF are ASCII codes for carriage return and line feed.
	CR = '\r' // used on Windows and some older macOS files
	LF = '\n' // used on Unix and macOS
)

var (
	// pre-computed lookup table for delimiters
	isSpaceDelimiter [256]bool
)

func init() {
	// initialize the lookup table for delimiters
	for _, ch := range []byte{CR, LF, ',', '(', ')', '\\', ':'} {
		isSpaceDelimiter[ch] = true
	}
}

type Parser struct {
	pos      int    // current position in the input (offset of next character to read)
	input    []byte // input data, not normalized and caller must not modify
	pushback []byte // pushback buffer for look-ahead input
}

func New(options ...Option_t) (*Parser, error) {
	p := &Parser{}
	for _, opt := range options {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *Parser) Free() {
	// TODO: free the parser
}

type Node any

// Parse returns the root node of the parse tree.
func (p *Parser) Parse() (Node, error) {
	started := time.Now()
	noUnits := 0
	defer func() {
		log.Printf("pigeon: parse read %d unit sections\n", noUnits)
		log.Printf("pigeon: parse completed in %v\n", time.Since(started))
	}()

	// highest level loop in the parser looks for unit headings.
	// every time we find one, we pass control to a lower level of the parser.
	// that level returns the number of lines parsed and the highest level error it encountered.
	for line := p.NextLine(); line != nil; line = p.NextLine() {
		// is this line a unit heading?
		if _, ok := acceptUnitId(line); !ok {
			continue
		}
		noUnits++
		//// is this line a turn number?
		//if turnLine := turns.ParseTurnLine(path, line); turnLine != nil {
		//	rpt.Turn = &Turn_t{
		//		No:    turnLine.No,
		//		Year:  turnLine.Year,
		//		Month: turnLine.Month,
		//		Error: turnLine.Error,
		//	}
		//	// log.Printf("import: report: turn number: %+v\n", *turnLine)
		//	continue
		//}
		//// is this line a tribe movement?
		//// is this line a tribe follows?
		//// is this line a tribe goes to?
		//// is this line a fleet movement?
		//// is this line a scouting report?
		//// is this line a tribe status?
	}

	return nil, nil
}

// Lines returns the entire set of lines in the input.
// It returns a copy of the input data, so the caller is free to modify it.
// Does not update the parser's position.
func (p *Parser) Lines() (lines [][]byte) {
	originalPos := p.pos
	p.pos = 0
	for line := p.NextLine(); line != nil; line = p.NextLine() {
		lines = append(lines, line)
	}
	p.pos = originalPos
	return lines
}

// NextLine returns the next line from the input.
// The slice returned is a copy of the input data, so the caller is free to modify it.
// The returned slice never includes trailing whitespace or the line separator.
//
// Returns nil only at end of input.
func (p *Parser) NextLine() (line []byte) {
	if p == nil || !(p.pos < len(p.input)) {
		return nil
	} else if p.pushback != nil {
		// return the pushback buffer
		line, p.pushback = p.pushback, nil
		return line
	} else if !(p.pos < len(p.input)) {
		return nil
	}

	for prevCharWasDelimiter := false; p.pos < len(p.input); p.pos++ {
		if p.input[p.pos] == ' ' || p.input[p.pos] == '\t' {
			// skip run of spaces and tabs
			for ; p.pos < len(p.input) && p.input[p.pos] == ' ' || p.input[p.pos] == '\t'; p.pos++ {
				// consume the space or tab
			}
			// check if the space is significant or not
			nextCharIsDelimiter := p.pos >= len(p.input) || isSpaceDelimiter[p.input[p.pos]]
			isSignificantSpace := !(prevCharWasDelimiter || nextCharIsDelimiter)
			if isSignificantSpace {
				// space is significant, so keep it
				line = append(line, ' ')
			}
			p.pos-- // adjust for the outer loop increment
			continue
		} else if 'A' <= p.input[p.pos] && p.input[p.pos] <= 'Z' {
			line = append(line, p.input[p.pos]-'A'+'a') // convert to lowercase
			prevCharWasDelimiter = false                // reset the delimiter state
			continue
		} else if p.input[p.pos] != LF && p.input[p.pos] != CR {
			// write the current character and update the delimiter state
			line = append(line, p.input[p.pos])
			prevCharWasDelimiter = isSpaceDelimiter[p.input[p.pos]]
			continue
		}
		// note that we're going to treat a run of CR as a single CR.
		// this isn't exactly correct, but it's close enough for our purposes.
		// and it lets us handle old macOS files.
		for p.pos < len(p.input) && p.input[p.pos] == CR {
			p.pos++ // consume the CR
		}
		if p.pos < len(p.input) && p.input[p.pos] == LF {
			p.pos++ // consume the LF
		}
		break
	}

	if len(line) == 0 {
		// we have to return a non-nil slice if the input is not empty.
		return []byte{}
	}

	return line
}

func (p *Parser) peekLine() (line []byte) {
	if p.pushback == nil {
		p.pushback = p.NextLine()
	}
	return p.pushback
}

func (p *Parser) parseUnit() {}

func (p *Parser) parseUnitHeader() {}

func (p *Parser) parseCurrentTurn() {}

func (p *Parser) parseTribeMovement() {}

func (p *Parser) parseTribeFollows() {}

func (p *Parser) parseTribeGoesTo() {}

func (p *Parser) parseFleetMovement() {}

func (p *Parser) parseScoutingReport() {}

func (p *Parser) parseScoutMovement() {}

func (p *Parser) parseUnitStatus() {}
