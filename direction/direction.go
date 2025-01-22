// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package direction

import (
	"encoding/json"
	"fmt"
)

// Direction_e is an enum for the direction
type Direction_e int

const (
	None Direction_e = iota
	North
	NorthEast
	SouthEast
	South
	SouthWest
	NorthWest
)
const (
	NumDirections = int(NorthWest) + 1
)

// Directions is a helper for iterating over the directions
var Directions = []Direction_e{
	North,
	NorthEast,
	SouthEast,
	South,
	SouthWest,
	NorthWest,
}

// MarshalJSON implements the json.Marshaler interface.
func (d Direction_e) MarshalJSON() ([]byte, error) {
	return json.Marshal(EnumToString[d])
}

// MarshalText implements the encoding.TextMarshaler interface.
// This is needed for marshalling the enum as map keys.
//
// Note that this is called by the json package, unlike the UnmarshalText function.
func (d Direction_e) MarshalText() (text []byte, err error) {
	return []byte(EnumToString[d]), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Direction_e) UnmarshalJSON(data []byte) error {
	var s string
	var ok bool
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if *d, ok = StringToEnum[s]; !ok {
		return fmt.Errorf("invalid Direction %q", s)
	}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// This is needed for unmarshalling the enum as map keys.
//
// Note that this is never called; it just changes the code path in UnmarshalJSON.
func (d Direction_e) UnmarshalText(text []byte) error {
	panic("!")
}

// String implements the fmt.Stringer interface.
func (d Direction_e) String() string {
	if str, ok := EnumToString[d]; ok {
		return str
	}
	return fmt.Sprintf("Direction(%d)", int(d))
}

var (
	// EnumToString is a helper map for marshalling the enum
	EnumToString = map[Direction_e]string{
		None:      "",
		North:     "N",
		NorthEast: "NE",
		SouthEast: "SE",
		South:     "S",
		SouthWest: "SW",
		NorthWest: "NW",
	}
	// StringToEnum is a helper map for unmarshalling the enum
	StringToEnum = map[string]Direction_e{
		"":   None,
		"N":  North,
		"NE": NorthEast,
		"SE": SouthEast,
		"S":  South,
		"SW": SouthWest,
		"NW": NorthWest,
	}

	LowercaseToEnum = map[string]Direction_e{
		"n":  North,
		"ne": NorthEast,
		"se": SouthEast,
		"s":  South,
		"sw": SouthWest,
		"nw": NorthWest,
	}
)

// column direction vectors defines the vectors used to determine the coordinates
// of the neighboring column based on the direction and the odd/even column
// property of the starting hex.
//
// NB: grids and hexes start at 1, 1 so "odd" and "even" are based on the hex coordinates.

var OddColumnVectors = map[Direction_e][2]int{
	North:     [2]int{+0, -1}, // ## 1306 -> ## 1305
	NorthEast: [2]int{+1, -1}, // ## 1306 -> ## 1405
	SouthEast: [2]int{+1, +0}, // ## 1306 -> ## 1406
	South:     [2]int{+0, +1}, // ## 1306 -> ## 1307
	SouthWest: [2]int{-1, +0}, // ## 1306 -> ## 1206
	NorthWest: [2]int{-1, -1}, // ## 1306 -> ## 1205
}

var EvenColumnVectors = map[Direction_e][2]int{
	North:     [2]int{+0, -1}, // ## 1206 -> ## 1205
	NorthEast: [2]int{+1, +0}, // ## 1206 -> ## 1306
	SouthEast: [2]int{+1, +1}, // ## 1206 -> ## 1307
	South:     [2]int{+0, +1}, // ## 1206 -> ## 1207
	SouthWest: [2]int{-1, +1}, // ## 1206 -> ## 1107
	NorthWest: [2]int{-1, +0}, // ## 1206 -> ## 1106
}

// Add moves in the given direction and returns the new row and column.
// It always moves a single hex and allows for moving between grids and wrapping around the big map.
func Add(row, col int, d Direction_e) (int, int) {
	if d == None {
		return row, col
	} else if col%2 == 0 { // even column
		return row + EvenColumnVectors[d][1], col + EvenColumnVectors[d][0]
	}
	// odd column
	return row + OddColumnVectors[d][1], col + OddColumnVectors[d][0]
}
