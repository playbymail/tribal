// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lemon

import (
	"log"
)

type Parser struct {
	source    string
	input     []byte
	state     *State
	tokenizer *Tokenizer
}

func NewParser(source string, input []byte, tokenizer *Tokenizer) *Parser {
	p := &Parser{source: source, tokenizer: tokenizer}

	// convert to lowercase, compress spaces, and fix line endings
	data := make([]byte, 0, len(input))
	for pos, prevCharWasDelimiter := 0, false; pos < len(input); pos++ {
		if input[pos] == ' ' || input[pos] == '\t' {
			// skip run of spaces and tabs
			for ; pos < len(input) && input[pos] == ' ' || input[pos] == '\t'; pos++ {
				// consume the space or tab
			}
			// the space is significant only if the previous character wasn't a delimiter
			// and the next one isn't, either.
			nextCharIsDelimiter := pos >= len(input) || isSpaceDelimiter[input[pos]]
			isSignificantSpace := !(prevCharWasDelimiter || nextCharIsDelimiter)
			if isSignificantSpace {
				// space is significant, so keep it
				data = append(data, ' ')
			}
			pos-- // adjust for the outer loop increment
			continue
		} else if 'A' <= input[pos] && input[pos] <= 'Z' {
			data = append(data, input[pos]-'A'+'a') // convert to lowercase
			prevCharWasDelimiter = false            // reset the delimiter state
			continue
		} else if input[pos] != LF && input[pos] != CR {
			// write the current character and update the delimiter state
			data = append(data, input[pos])
			prevCharWasDelimiter = isSpaceDelimiter[input[pos]]
			continue
		}
		// note that we're going to treat a run of CR as a single CR.
		// this isn't exactly correct, but it's close enough for our purposes.
		// and it lets us handle old macOS files.
		for pos < len(input) && input[pos] == CR {
			pos++ // consume the CR
		}
		if pos < len(input) && input[pos] == LF {
			pos++ // consume the LF
		}
	}

	return p
}

// normalize converts the input to lowercase, compress spaces, and fix line endings.
// It also removes blank lines and lines starting with a space.
func normalize(input []byte) []byte {
	// data is a temporary buffer that we use to build the normalized input.
	data := make([]byte, 0, len(input))

	// pretend that we have a newline at the beginning of the input.
	previousChar := byte(LF)

	for pos, prevCharWasDelimiter := 0, false; pos < len(input); pos++ {
		// if we're at the start of a line (previous char was LF) and see space/tab,
		// skip the entire line
		if previousChar == LF && (input[pos] == ' ' || input[pos] == '\t') {
			// skip until we find LF or reach end of input
			for pos < len(input) && input[pos] != LF {
				pos++
			}
			// stay at the LF so the main loop's increment will move past it
			pos--
			previousChar = input[pos]
			continue
		}
		if input[pos] == ' ' || input[pos] == '\t' {
			// skip run of spaces and tabs
			for ; pos < len(input) && input[pos] == ' ' || input[pos] == '\t'; pos++ {
				// consume the space or tab
			}
			// the space is significant only if the previous character wasn't a delimiter
			// and the next one isn't, either.
			nextCharIsDelimiter := pos >= len(input) || isSpaceDelimiter[input[pos]]
			isSignificantSpace := !(prevCharWasDelimiter || nextCharIsDelimiter)
			if isSignificantSpace {
				// space is significant, so keep it
				data = append(data, ' ')
			}
			pos-- // adjust for the outer loop increment
			previousChar = input[pos]
			continue
		} else if 'A' <= input[pos] && input[pos] <= 'Z' {
			data = append(data, input[pos]-'A'+'a') // convert to lowercase
			prevCharWasDelimiter = false            // reset the delimiter state
			previousChar = input[pos]
			continue
		} else if input[pos] != LF && input[pos] != CR {
			// write the current character and update the delimiter state
			data = append(data, input[pos])
			prevCharWasDelimiter = isSpaceDelimiter[input[pos]]
			previousChar = input[pos]
			continue
		}

		previousChar = input[pos]

		// note that we're going to treat a run of CR as a single CR.
		// this isn't exactly correct, but it's close enough for our purposes.
		// and it lets us handle old macOS files.
		for pos < len(input) && input[pos] == CR {
			pos++ // consume the CR
		}
		if pos < len(input) && input[pos] == LF {
			pos++ // consume the LF
		}
	}

	return data
}

// Parse will panic if it needs to report an error.
// Otherwise, it will return the root node of the parse tree.
func (p *Parser) Parse() Node {
	log.Printf("joy\n")
	p.parseReport()
	return Node{}
}

func (p *Parser) WriteBuffer(path string) error {
	return nil
}

type report struct {
	sections []*unit_section
}

func (p *Parser) parseReport() *report {
	sections := p.parseUnitSections()
	return &report{
		sections: sections,
	}
}

func (p *Parser) parseUnitSections() (sections []*unit_section) {
	for {
		p.tokenizer.runToEndOfLine()
		if p.tokenizer.peekch() != LF {
			break
		}
		section := p.parseUnitSection()
		if section == nil {
			break
		}
		sections = append(sections, section)
	}
	return sections
}

type unit_section struct{}

func (p *Parser) parseUnitSection() *unit_section {
	if p.tokenizer.peekch() != LF {
		return nil
	}
	p.tokenizer.getch() // consume LF
	return &unit_section{}
}

type unit_header struct{}

func (p *Parser) parseUnitHeader() unit_header {
	return unit_header{}
}

type unit_header_id struct{}

func (p *Parser) parseUnitHeaderId() unit_header_id {
	return unit_header_id{}
}

type unit_name struct{}

func (p *Parser) parseUnitName() unit_name {
	return unit_name{}
}

type current_hex struct{}

func (p *Parser) parseCurrentHex() current_hex {
	return current_hex{}
}

type previous_hex struct{}

func (p *Parser) parsePreviousHex() previous_hex {
	return previous_hex{}
}

type unit_status struct{}

func (p *Parser) parseUnitStatus() unit_status {
	return unit_status{}
}
