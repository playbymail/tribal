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
func ParseTribeGoesTo(id ast.UnitId_t, from, to ast.Coordinates_t, input []byte) (*ast.GoesTo_t, error) {
	if match := reUnitGoesTo.FindSubmatch(input); match == nil {
		return nil, ast.ErrInvalidUnitGoesTo
	} else if coords, err := ast.TextToCoordinates(match[1]); err != nil {
		return nil, err
	} else {
		return &ast.GoesTo_t{
			Id:     id,
			From:   from,
			GoesTo: coords,
			To:     to,
		}, nil
	}
}
