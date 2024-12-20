// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package adapters

import (
	"bytes"
	"github.com/playbymail/tribal/parser/ast"
)

// TextToCoordinates converts text to coordinates.
// Note that grid, row, and column are 1-based, not 0-based.
// We return an error if the input is invalid.
func TextToCoordinates(text []byte) (ast.Coordinates_t, error) {
	if bytes.Equal(text, []byte{'n', '/', 'a'}) {
		return ast.Coordinates_t{}, nil
	} else if len(text) != 7 {
		return ast.Coordinates_t{}, ast.ErrInvalidCoordinates
	}
	c := ast.Coordinates_t{
		GridRow:    int(text[0]-'a') + 1,
		GridColumn: int(text[1]-'a') + 1,
		Column:     int(text[3]-'0')*10 + int(text[4]-'0'),
		Row:        int(text[5]-'0')*10 + int(text[6]-'0'),
	}
	// obscured grid gets zero for row and column
	if bytes.HasPrefix(text, []byte{'#', '#'}) {
		c.GridRow, c.GridColumn = 0, 0
	} else if !(0 <= c.GridRow && c.GridRow <= 26) {
		return c, ast.ErrInvalidCoordinates
	} else if !(0 <= c.GridColumn && c.GridColumn <= 26) {
		return c, ast.ErrInvalidCoordinates
	}
	// all grids must have valid columns and rows
	if !(1 <= c.Column && c.Column <= 30) {
		return c, ast.ErrInvalidCoordinates
	} else if !(1 <= c.Row && c.Row <= 20) {
		return c, ast.ErrInvalidCoordinates
	}
	return c, nil
}
