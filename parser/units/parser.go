// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package units implements the parser for the units section of a turn report.
package units

import (
	"bytes"
	"fmt"
	"github.com/playbymail/tribal/parser/ast"
	"log"
)

//go:generate pigeon -o grammar.go grammar.peg

// ParseUnitHeading parses a line of text and returns a UnitHeading_t if the line is a unit heading.
// We do a quick check to see if the line starts with a unit declaration.
func ParseUnitHeading(path string, input []byte) *ast.UnitHeading_t {
	id, ok := acceptUnitId(input)
	if !ok {
		//log.Printf("puh: %s: %v\n", input, false)
		return nil
	}
	//log.Printf("puh: %s: unit %q\n", input, id)

	// we have a unit heading. we will return it even if there are errors parsing the coordinates.
	unitHeading := &ast.UnitHeading_t{Id: id}

	// parse the line, capturing any errors
	if v, err := Parse(path, input); err != nil {
		unitHeading.Error = err
		log.Printf("puh: %s: parse error %v\n", input, err)
		return unitHeading
	} else if uh, ok := v.(*ast.UnitHeading_t); !ok {
		panic(fmt.Sprintf("assert(%T == *UnitHeading_t)", v))
	} else {
		unitHeading.Name = uh.Name
		unitHeading.CurrentHex = uh.CurrentHex
		unitHeading.PreviousHex = uh.PreviousHex
		unitHeading.Error = uh.Error
	}

	return unitHeading
}

// acceptUnitId returns the unit id and true if the line starts with
// a unit type followed by a unit id and a comma.
func acceptUnitId(input []byte) (ast.UnitId_t, bool) {
	if va, err := Parse("unit_id", input, Entrypoint("UnitDeclaration")); err != nil {
		// ignore the error. we know this means that a declaration was not found.
		return ast.UnitId_t(""), false
	} else if unitId, ok := va.(ast.UnitId_t); !ok {
		panic(fmt.Sprintf("assert(%T == UnitId_t)", va))
	} else {
		return unitId, true
	}
}

// tokenToCoords converts a token to a coordinate.
// it avoids error checks because it should only be called from the parser;
// we can assume that only valid tokens are passed in.
// note that grid, row, and column are 1-based, not 0-based.
// note that we do not check for out of bounds coordinates.
func tokenToCoords(token []byte) (ast.Coords_t, error) {
	//log.Printf("tokenToCoords: %d %q\n", len(token), token)
	c := ast.Coords_t{
		GridRow:    int(token[0]-'a') + 1,
		GridColumn: int(token[1]-'a') + 1,
		Column:     int(token[3]-'0')*10 + int(token[4]-'0'),
		Row:        int(token[5]-'0')*10 + int(token[6]-'0'),
	}
	// obscured grid gets zero for row and column
	if bytes.HasPrefix(token, []byte{'#', '#'}) {
		c.GridRow, c.GridColumn = 0, 0
	}
	if !(0 <= c.GridRow && c.GridRow <= 26) {
		return c, ast.ErrInvalidCoordinates
	} else if !(0 <= c.GridColumn && c.GridColumn <= 26) {
		return c, ast.ErrInvalidCoordinates
	} else if !(1 <= c.Column && c.Column <= 30) {
		return c, ast.ErrInvalidCoordinates
	} else if !(1 <= c.Row && c.Row <= 20) {
		return c, ast.ErrInvalidCoordinates
	}
	return c, nil
}
