// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package ast

type Turn_t struct {
	No    int
	Year  int
	Month int
	Error error // highest level error encountered while parsing the turn
}

type UnitSection_t struct {
	Unit  *UnitHeading_t // unit header from the parser
	Error error          // highest level error encountered while parsing the unit
}

type UnitHeading_t struct {
	Id          UnitId_t
	Name        UnitName_t     // optional name field
	PreviousHex *PreviousHex_t // location of unit at the beginning of the turn
	CurrentHex  *CurrentHex_t  // location of unit at the end of the turn
	Error       error          // highest level error encountered while parsing the unit
}

type ClanId_t string
type UnitId_t string
type UnitName_t string

type Coords_t struct {
	GridRow    int   // 1-based, A ... Z -> 1 ... 26
	GridColumn int   // 1-based, A ... Z -> 1 ... 26
	Column     int   // 1-based, 1 ... 30
	Row        int   // 1-based, 1 ... 21
	Error      error // highest level error encountered while parsing the coordinates
}

type CurrentHex_t struct {
	IsNA       bool
	IsObscured bool
	Coords     Coords_t
	Error      error // highest level error encountered while parsing the coordinates
}

type PreviousHex_t struct {
	IsNA       bool
	IsObscured bool
	Coords     Coords_t
	Error      error // highest level error encountered while parsing the coordinates
}
