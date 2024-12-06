// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package turns

import (
	"bytes"
	"fmt"
)

//go:generate pigeon -o grammar.go grammar.peg

type TurnLine_t struct {
	No    int
	Year  int
	Month int
	Error error
}

// ParseTurnLine parses a line of text and returns a TurnLine_t if the line is a turn line.
func ParseTurnLine(path string, input []byte) *TurnLine_t {
	// cheap check to see if the input is a turn line
	if len(input) == 0 || !bytes.HasPrefix(input, []byte("current turn ")) {
		return nil
	}
	iTurnLine, err := Parse(path, input)
	if err != nil {
		// silently ignore errors since we only care about matches
		return nil
	}
	turnLine, ok := iTurnLine.(*TurnLine_t)
	if !ok {
		panic(fmt.Sprintf("assert(%T == *TurnLine_t)", iTurnLine))
	}
	return turnLine
}
