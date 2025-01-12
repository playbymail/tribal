// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lemon

import (
	"bytes"
	"github.com/playbymail/tribal/is"
	"github.com/playbymail/tribal/norm"
	"os"
)

const (
	// CR and LF are ASCII codes for carriage return and line feed.
	CR = '\r' // used on Windows and some older macOS files
	LF = '\n' // used on Unix and macOS
)

type Tokenizer struct {
	pos      int
	data     []byte
	pushback []byte
}

// NewTokenizer returns a new tokenizer for the given input.
// A CR is pushed onto the pushback buffer so that the first read
// will think that it is starting with a newline.
func NewTokenizer(input []byte) *Tokenizer {
	input = norm.NormalizeSpaces(input)
	input = norm.NormalizeCase(input)
	input = norm.LineEndings(input)
	var lines [][]byte
	for _, line := range norm.RemoveEmptyLines(bytes.Split(input, []byte{LF})) {
		if is.FleetMovement(line) {
			lines = append(lines, norm.FleetMovement(line))
		} else if is.ScoutLine(line) {
			lines = append(lines, norm.ScoutMovement(line))
		} else if is.TribeMovement(line) {
			lines = append(lines, norm.TribeMovement(line))
		} else if is.TribeFollows(line) {
			lines = append(lines, line)
		} else if is.TribeGoesTo(line) {
			lines = append(lines, line)
		} else if is.TurnHeader(line) {
			lines = append(lines, line)
		} else if is.UnitHeader(line) {
			lines = append(lines, line)
		} else if is.UnitStatus(line) {
			lines = append(lines, norm.UnitStatus(line))
		}
	}
	input = bytes.Join(lines, []byte{LF})

	return &Tokenizer{
		data:     input,
		pushback: []byte{LF},
	}
}

// getch returns the next character in the input.
// It returns 0 if there are no more characters.
func (t *Tokenizer) getch() (ch byte) {
	if len(t.pushback) != 0 {
		ch, t.pushback = t.pushback[len(t.pushback)-1], t.pushback[:len(t.pushback)-1]
		return ch
	} else if t.pos < len(t.data) {
		ch, t.pos = t.data[t.pos], t.pos+1
		return ch
	}
	return 0
}

// peekch returns the next character in the input.
// It returns 0 if there are no more characters.
func (t *Tokenizer) peekch() byte {
	if len(t.pushback) != 0 {
		return t.pushback[len(t.pushback)-1]
	} else if t.pos < len(t.data) {
		return t.data[t.pos]
	}
	return 0
}

// ungetch pushes a character back onto the input.
func (t *Tokenizer) ungetch(ch byte) {
	t.pushback = append(t.pushback, ch)
}

// runToEndOfLine discards all input up to the next newline or the end of the input.
// Leaves the input pointing at the newline or at the end of the input.
func (t *Tokenizer) runToEndOfLine() {
	for ch := t.getch(); ch != 0; ch = t.getch() {
		if ch == LF { // found a newline; leave the input pointing at it
			t.ungetch(ch)
			return
		}
	}
	return
}

func (t *Tokenizer) WriteTo(path string) error {
	return os.WriteFile(path, t.data, 0644)
}

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
