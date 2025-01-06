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

type Unit_t struct {
	Kind UnitKind_t
	No   int // range 1 ... 9999
	Seq  int // range 0 ... 9
}

type UnitKind_t int

const (
	Tribe UnitKind_t = iota
	Courier
	Element
	Fleet
	Garrison
)

type ClanId_t string
type UnitId_t string
type UnitName_t string

type Coordinates_t struct {
	GridRow    int // 1-based, A ... Z -> 1 ... 26
	GridColumn int // 1-based, A ... Z -> 1 ... 26
	Column     int // 1-based, 1 ... 30
	Row        int // 1-based, 1 ... 21
}

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
