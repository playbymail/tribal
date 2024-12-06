// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package parser

import (
	"bytes"
	"fmt"
)

type Report_t struct {
	Hash     string // SHA1 hash of the report
	Turn     *Turn_t
	Units    []*Unit_t
	Sections []*Section_t
	Error    error    // highest level error encountered while parsing the report
	Lines    [][]byte // copy of the input after normalization
	options  Options_t
}

type Turn_t struct {
	No    int
	Year  int
	Month int
	Error error // highest level error encountered while parsing the turn
}

func (t *Turn_t) String() string {
	if t == nil {
		return (&Turn_t{}).String()
	}
	return fmt.Sprintf("%04d-%02d", t.Year, t.Month)
}

type Unit_t struct {
	Id         UnitId_t
	CurrentHex *Coords_t // location of the unit at the end of the turn
	PriorHex   *Coords_t // location of the unit at the beginning of the turn, if known
	Turn       *Turn_t   // nil unless the parser finds a turn number in the section
	Children   []*Unit_t // child elements of this unit
	Error      error     // highest level error encountered while parsing the unit
}

type Section_t struct {
	Unit  UnitId_t
	Error error
}

type Coords_t struct {
	GridRow    int // 1-based, A ... Z -> 1 ... 26
	GridColumn int // 1-based, A ... Z -> 1 ... 26
	Column     int // 1-based, 1 ... 30
	Row        int // 1-based, 1 ... 21
}

// tokenToCoords converts a token to a coordinate.
// it avoids error checks because it should only be called from the parser;
// we can assume that only valid tokens are passed in.
// note that grid, row, and column are 1-based, not 0-based.
func tokenToCoords(token []byte) Coords_t {
	if len(token) != 7 {
		return Coords_t{}
	}
	// grid can be "##" or "aa" .. "zz"
	var gridRow, gridColumn int
	if bytes.HasPrefix(token, []byte{'#', '#'}) {
		// obscured grid gets zero for row and column
		gridRow, gridColumn = 0, 0
	} else {
		// visible grid. id will be aa .. zz
		gridRow = int(token[0]-'a') + 1
		gridColumn = int(token[1]-'z') + 1
	}
	// remainder of token is column and row.
	// check bounds after conversion and return a zero location if out of bounds.
	column := int(token[3]-'0')*10 + int(token[4]-'0')
	if !(1 <= column && column <= 30) {
		return Coords_t{}
	}
	row := int(token[5]-'0')*10 + int(token[6]-'0')
	if !(1 <= row && row <= 21) {
		return Coords_t{}
	}
	return Coords_t{GridRow: gridRow, GridColumn: gridColumn, Column: column, Row: row}
}

// IsNA returns true if the location is not set.
// This is the same as IsZero.
func (c Coords_t) IsNA() bool {
	return c == Coords_t{}
}

// IsZero returns true if the location is not set.
func (c Coords_t) IsZero() bool {
	return c == Coords_t{}
}

func (c Coords_t) String() string {
	if c.IsNA() {
		return "n/a"
	} else if c.GridRow == 0 && c.GridColumn == 0 {
		return fmt.Sprintf("## %02d%02d", c.Column, c.Row)
	}
	return fmt.Sprintf("%c%c %02d%02d", 'A'+c.GridRow-1, 'A'+c.GridColumn-1, c.Column, c.Row)
}

type SectionHeader_t struct {
	UnitId     UnitId_t
	CurrentHex Coords_t
}

type UnitId_t string
