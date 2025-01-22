// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/domains"
)

// Coordinates_t defines a location (grid, row, and column) of a tile in a turn report
type Coordinates_t struct {
	GridRow    int // 1-based, A ... Z -> 1 ... 26
	GridColumn int // 1-based, A ... Z -> 1 ... 26
	Column     int // 1-based, 1 ... 30
	Row        int // 1-based, 1 ... 21
}

func (c Coordinates_t) IsValidGrid() bool {
	return c.GridRow != 0 && c.GridColumn != 0
}

func (c Coordinates_t) IsZero() bool {
	return c.GridRow == 0 && c.GridColumn == 0 && c.Column == 0 && c.Row == 0
}

func (c Coordinates_t) String() string {
	if c.IsZero() {
		return "n/a"
	} else if c.GridRow == 0 && c.GridColumn == 0 {
		return fmt.Sprintf("## %02d%02d", c.Column, c.Row)
	}
	return fmt.Sprintf("%c%c %02d%02d", c.GridRow+'A'-1, c.GridColumn+'A'-1, c.Column, c.Row)
}

// TextToCoordinates converts text to coordinates.
// Note that grid, row, and column are 1-based, not 0-based.
// We return an error if the input is invalid.
func TextToCoordinates(text []byte) (Coordinates_t, error) {
	if bytes.Equal(text, []byte{'n', '/', 'a'}) {
		return Coordinates_t{}, nil
	} else if len(text) != 7 {
		return Coordinates_t{}, domains.ErrInvalidCoordinates
	}
	c := Coordinates_t{
		GridRow:    int(text[0]-'a') + 1,
		GridColumn: int(text[1]-'a') + 1,
		Column:     int(text[3]-'0')*10 + int(text[4]-'0'),
		Row:        int(text[5]-'0')*10 + int(text[6]-'0'),
	}
	// obscured grid gets zero for row and column
	if bytes.HasPrefix(text, []byte{'#', '#'}) {
		c.GridRow, c.GridColumn = 0, 0
	} else if !(0 <= c.GridRow && c.GridRow <= 26) {
		return c, domains.ErrInvalidCoordinates
	} else if !(0 <= c.GridColumn && c.GridColumn <= 26) {
		return c, domains.ErrInvalidCoordinates
	}
	// all grids must have valid columns and rows
	if !(1 <= c.Column && c.Column <= 30) {
		return c, domains.ErrInvalidCoordinates
	} else if !(1 <= c.Row && c.Row <= 20) {
		return c, domains.ErrInvalidCoordinates
	}
	return c, nil
}

func (c Coordinates_t) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c Coordinates_t) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

//func (c *Coordinates_t) UnmarshalJSON(data []byte) error {
//	return json.Unmarshal(data, &m.Field)
//}

func (c Coordinates_t) Move(d direction.Direction_e) Coordinates_t {
	if d == direction.None {
		return c
	} else if c.IsZero() {
		return c
	}

	invalidGrid := !c.IsValidGrid()
	gcol, grow := c.GridColumn, c.GridRow
	row, col := direction.Add(c.Row, c.Column, d)

	// allow for moving between grids and wrapping around the edges of the big map
	if col < 1 {
		col = 30
		if invalidGrid {
			// invalid grid, so don't bother to check for wrapping
		} else if gcol == 1 {
			gcol = 26
		} else {
			gcol = gcol - 1
		}
	} else if col > 30 {
		col = 1
		if invalidGrid {
			// invalid grid, so don't bother to check for wrapping
		} else if gcol == 26 {
			gcol = 1
		} else {
			gcol = gcol + 1
		}
	}
	if row < 1 {
		row = 21
		if invalidGrid {
			// invalid grid, so don't bother to check for wrapping
		} else if grow == 1 {
			grow = 26
		} else {
			grow = grow - 1
		}
	} else if row > 21 {
		row = 1
		if invalidGrid {
			// invalid grid, so don't bother to check for wrapping
		} else if grow == 26 {
			grow = 1
		} else {
			grow = grow + 1
		}
	}

	return Coordinates_t{GridRow: grow, GridColumn: gcol, Column: col, Row: row}
}
