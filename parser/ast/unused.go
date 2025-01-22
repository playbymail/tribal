// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ast

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

type UnitKind_t int

const (
	Tribe UnitKind_t = iota
	Courier
	Element
	Fleet
	Garrison
)

type CurrentHex_t struct {
	IsNA       bool
	IsObscured bool
	Coords     Coordinates_t
	Error      error // highest level error encountered while parsing the coordinates
}

type PreviousHex_t struct {
	IsNA       bool
	IsObscured bool
	Coords     Coordinates_t
	Error      error // highest level error encountered while parsing the coordinates
}
