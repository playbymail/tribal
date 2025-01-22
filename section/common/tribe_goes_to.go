// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package common

import (
	"github.com/playbymail/tribal/parser/ast"
	"regexp"
)

var (
	// reUnitGoesTo is the regular expression for a unit goes to line.
	reUnitGoesTo = regexp.MustCompile(`^tribe goes to ([a-z]{2} \d{4})$`)
)

// ParseTribeGoesTo parses the tribe goes to line.
//
//	"tribe goes to" Coordinates
func ParseTribeGoesTo(path string, input []byte) (*ast.Coordinates_t, error) {
	if match := reUnitGoesTo.FindSubmatch(input); match == nil {
		return nil, ast.ErrInvalidUnitGoesTo
	} else if coords, err := ast.TextToCoordinates(match[1]); err != nil {
		return nil, err
	} else {
		return &coords, err
	}
}
