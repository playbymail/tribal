// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package units implements the parser for the units section of a turn report.
package units

import (
	"fmt"
	"github.com/playbymail/tribal/parser/ast"
	"log"
)

//go:generate pigeon -o grammar.go grammar.peg

// ParseUnitHeading parses a line of text and returns a UnitHeading_t if the line is a unit heading.
// We do a quick check to see if the line starts with a unit declaration.
func ParseUnitHeading(path string, input []byte) *ast.UnitHeading_t {
	log.Printf("puh: %q\n", input)
	id, ok := acceptUnitId(input)
	if !ok {
		log.Printf("puh: %q: %v\n", input, false)
		return nil
	}
	log.Printf("puh: %q: unit %q\n", input, id)

	// we have a unit heading. we will return it even if there are errors parsing the coordinates.
	unitHeading := &ast.UnitHeading_t{Id: id}

	// parse the line, capturing any errors
	if v, err := Parse(path, input); err != nil {
		unitHeading.Error = err
		log.Printf("puh: %q: parse error %v\n", input, err)
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
