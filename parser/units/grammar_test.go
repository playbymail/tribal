// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package units_test

import (
	"errors"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/parser/units"
	"testing"
)

// implements tests for parsing a unit section

func TestUnitHeaderParser(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected *ast.UnitHeading_t
		err      error
	}{
		{
			name:  "tribe",
			input: `tribe 0987,,current hex = ab 1013,(previous hex = ab 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "0987",
				PreviousHex: &ast.PreviousHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 12, Row: 11}},
				CurrentHex:  &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
		{
			name:  "courier",
			input: `courier 0987c1,,current hex = ab 1013,(previous hex = ab 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "0987c1",
				PreviousHex: &ast.PreviousHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 12, Row: 11}},
				CurrentHex:  &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
		{
			name:  "element",
			input: `element 0987e1,,current hex = ab 1013,(previous hex = ab 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "0987e1",
				PreviousHex: &ast.PreviousHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 12, Row: 11}},
				CurrentHex:  &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
		{
			name:  "fleet",
			input: `fleet 0987f1,,current hex = ab 1013,(previous hex = ab 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "0987f1",
				PreviousHex: &ast.PreviousHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 12, Row: 11}},
				CurrentHex:  &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
		{
			name:  "garrison",
			input: `garrison 0987g1,,current hex = ab 1013,(previous hex = ab 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "0987g1",
				PreviousHex: &ast.PreviousHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 12, Row: 11}},
				CurrentHex:  &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
		{
			name:  "n/a coords",
			input: `garrison 0987g2,,current hex = ab 1013,(previous hex = n/a)`,
			expected: &ast.UnitHeading_t{
				Id:          "0987g2",
				PreviousHex: &ast.PreviousHex_t{IsNA: true},
				CurrentHex:  &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
		{
			name:  "obscured grid",
			input: `tribe 1987,,current hex = ab 1013,(previous hex = ## 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "1987",
				PreviousHex: &ast.PreviousHex_t{IsObscured: true, Coords: ast.Coords_t{Column: 12, Row: 11}},
				CurrentHex:  &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
		{
			name:  "obscured grids",
			input: `tribe 1987,,current hex = ## 1013,(previous hex = ## 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "1987",
				PreviousHex: &ast.PreviousHex_t{IsObscured: true, Coords: ast.Coords_t{Column: 12, Row: 11}},
				CurrentHex:  &ast.CurrentHex_t{IsObscured: true, Coords: ast.Coords_t{Column: 10, Row: 13}},
			},
		},
		{
			name:  "invalid current hex",
			input: `tribe 1987,,current hex = na,(previous hex = ## 1211)`,
			expected: &ast.UnitHeading_t{
				Id:          "1987",
				PreviousHex: &ast.PreviousHex_t{IsObscured: true, Coords: ast.Coords_t{Column: 12, Row: 11}},
			},
		},
		{
			name:  "invalid previous hex",
			input: `tribe 1987,,current hex = ab 1013,(previous hex = ab 11)`,
			expected: &ast.UnitHeading_t{
				Id:         "1987",
				CurrentHex: &ast.CurrentHex_t{Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 10, Row: 13}},
			},
		},
	} {
		rslt, err := units.ParseUnitHeading(tc.name, []byte(tc.input))

		// errors are fragile
		if err != nil || tc.err != nil {
			if tc.err != nil && err != nil {
				t.Errorf("%s: error: expected %v, got %v", tc.name, tc.err, err)
			} else if tc.err != nil {
				t.Errorf("%s: error: expected %v, got nil", tc.name, tc.err)
			} else if err != nil {
				t.Errorf("%s: error: expected nil, got %v", tc.name, err)
			}
			continue
		}

		if rslt == nil || tc.expected == nil {
			if rslt == nil && tc.expected == nil {
				// okay
			} else if rslt == nil {
				t.Errorf("%s: expected result, got nil", tc.name)
			} else if tc.expected == nil {
				t.Errorf("%s: expected nil, got %+v", tc.name, *rslt)
			}
			continue
		}

		if tc.expected.Id != rslt.Id {
			t.Errorf("%s: unit_id: expected %q, got %q", tc.name, tc.expected.Id, rslt.Id)
		}

		if tc.expected.CurrentHex == nil && rslt.CurrentHex == nil {
			// ok
		} else if tc.expected.CurrentHex == nil {
			t.Errorf("%s: current_hex: expected nil, got %+v", tc.name, *rslt.CurrentHex)
		} else if rslt.CurrentHex == nil {
			t.Errorf("%s: current_hex: expected current hex, got nil", tc.name)
		} else if *rslt.CurrentHex != *tc.expected.CurrentHex {
			if rslt.CurrentHex.IsNA != tc.expected.CurrentHex.IsNA {
				t.Errorf("%s: current_hex: n/a: expected %v, got %v", tc.name, tc.expected.CurrentHex.IsNA, rslt.CurrentHex.IsNA)
			}
			if rslt.CurrentHex.IsObscured != tc.expected.CurrentHex.IsObscured {
				t.Errorf("%s: current_hex: obscured: expected %v, got %v", tc.name, tc.expected.CurrentHex.IsObscured, rslt.CurrentHex.IsObscured)
			}
			if rslt.CurrentHex.Coords.GridRow != tc.expected.CurrentHex.Coords.GridRow {
				t.Errorf("%s: current_hex: grid_row: expected %d, got %d", tc.name, tc.expected.CurrentHex.Coords.GridRow, rslt.CurrentHex.Coords.GridRow)
			}
			if rslt.CurrentHex.Coords.GridColumn != tc.expected.CurrentHex.Coords.GridColumn {
				t.Errorf("%s: current_hex: grid_col: expected %d, got %d", tc.name, tc.expected.CurrentHex.Coords.GridColumn, rslt.CurrentHex.Coords.GridColumn)
			}
			if rslt.CurrentHex.Coords.Column != tc.expected.CurrentHex.Coords.Column {
				t.Errorf("%s: current_hex: map_col: expected %d, got %d", tc.name, tc.expected.CurrentHex.Coords.Column, rslt.CurrentHex.Coords.Column)
			}
			if rslt.CurrentHex.Coords.Row != tc.expected.CurrentHex.Coords.Row {
				t.Errorf("%s: current_hex: map_row: expected %d, got %d", tc.name, tc.expected.CurrentHex.Coords.Row, rslt.CurrentHex.Coords.Row)
			}
		}

		if tc.expected.PreviousHex == nil && rslt.PreviousHex == nil {
			// ok
		} else if tc.expected.PreviousHex == nil {
			t.Errorf("%s: previous_hex: expected nil, got %+v", tc.name, *rslt.PreviousHex)
		} else if rslt.PreviousHex == nil {
			t.Errorf("%s: previous_hex: expected previous hex, got nil", tc.name)
		} else if *rslt.PreviousHex != *tc.expected.PreviousHex {
			if rslt.PreviousHex.IsNA != tc.expected.PreviousHex.IsNA {
				t.Errorf("%s: previous_hex: n/a: expected %v, got %v", tc.name, tc.expected.PreviousHex.IsNA, rslt.PreviousHex.IsNA)
			}
			if rslt.PreviousHex.IsObscured != tc.expected.PreviousHex.IsObscured {
				t.Errorf("%s: previous_hex: obscured: expected %v, got %v", tc.name, tc.expected.PreviousHex.IsObscured, rslt.PreviousHex.IsObscured)
			}
			if rslt.PreviousHex.Coords.GridRow != tc.expected.PreviousHex.Coords.GridRow {
				t.Errorf("%s: previous_hex: grid_row: expected %d, got %d", tc.name, tc.expected.PreviousHex.Coords.GridRow, rslt.PreviousHex.Coords.GridRow)
			}
			if rslt.PreviousHex.Coords.GridColumn != tc.expected.PreviousHex.Coords.GridColumn {
				t.Errorf("%s: previous_hex: grid_col: expected %d, got %d", tc.name, tc.expected.PreviousHex.Coords.GridColumn, rslt.PreviousHex.Coords.GridColumn)
			}
			if rslt.PreviousHex.Coords.Column != tc.expected.PreviousHex.Coords.Column {
				t.Errorf("%s: previous_hex: map_col: expected %d, got %d", tc.name, tc.expected.PreviousHex.Coords.Column, rslt.PreviousHex.Coords.Column)
			}
			if rslt.PreviousHex.Coords.Row != tc.expected.PreviousHex.Coords.Row {
				t.Errorf("%s: previous_hex: map_row: expected %d, got %d", tc.name, tc.expected.PreviousHex.Coords.Row, rslt.PreviousHex.Coords.Row)
			}
		}

	}
}

func TestCurrentHexParser(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected ast.CurrentHex_t
		err      error
	}{
		{
			name:  "grid exposed",
			input: `current hex = ab 1211`,
			expected: ast.CurrentHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 12, Row: 11},
			},
		},
		{
			name:  "grid obscured",
			input: `current hex = ## 1211`,
			expected: ast.CurrentHex_t{
				IsObscured: true,
				Coords:     ast.Coords_t{GridRow: 0, GridColumn: 0, Column: 12, Row: 11},
			},
		},
		{
			name:  "n/a",
			input: `current hex = n/a`,
			expected: ast.CurrentHex_t{
				IsNA: true,
			},
		},
		{
			name:     "invalid grid row",
			input:    `current hex = 9a 1211`,
			expected: ast.CurrentHex_t{},
			err:      ast.ErrNoMatch,
		},
		{
			name:     "invalid grid column",
			input:    `current hex = a9 1211`,
			expected: ast.CurrentHex_t{},
			err:      ast.ErrNoMatch,
		},
		{
			name:  "invalid column low",
			input: `current hex = aa 0011`,
			expected: ast.CurrentHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 0, Row: 11},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
		{
			name:  "invalid column high",
			input: `current hex = aa 3111`,
			expected: ast.CurrentHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 31, Row: 11},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
		{
			name:  "invalid row low",
			input: `current hex = aa 1100`,
			expected: ast.CurrentHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 11, Row: 0},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
		{
			name:  "invalid row high",
			input: `current hex = aa 1122`,
			expected: ast.CurrentHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 11, Row: 22},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
	} {
		var rslt ast.CurrentHex_t
		var ok bool
		if v, err := units.Parse(tc.name, []byte(tc.input), units.Entrypoint("CurrentHex")); err != nil {
			// errors are fragile
			if tc.err == nil {
				t.Errorf("%s: error: expected nil, got %v", tc.name, err)
			} else if !errors.Is(tc.err, ast.ErrNoMatch) {
				t.Errorf("%s: error: expected %v, got %v", tc.name, tc.err, err)
			}
			continue
		} else if tc.err != nil {
			t.Errorf("%s: error: expected %v, got nil", tc.name, tc.err)
			continue
		} else if rslt, ok = v.(ast.CurrentHex_t); !ok {
			t.Errorf("%s: expected ast.CurrentHex_t, got %T", tc.name, v)
			continue
		}
		if rslt != tc.expected {
			if rslt.IsNA != tc.expected.IsNA {
				t.Errorf("%s: n/a: expected %v, got %v", tc.name, tc.expected.IsNA, rslt.IsNA)
			}
			if rslt.IsObscured != tc.expected.IsObscured {
				t.Errorf("%s: obscured: expected %v, got %v", tc.name, tc.expected.IsObscured, rslt.IsObscured)
			}
			if rslt.Coords.GridRow != tc.expected.Coords.GridRow {
				t.Errorf("%s: grid_row: expected %d, got %d", tc.name, tc.expected.Coords.GridRow, rslt.Coords.GridRow)
			}
			if rslt.Coords.GridColumn != tc.expected.Coords.GridColumn {
				t.Errorf("%s: grid_col: expected %d, got %d", tc.name, tc.expected.Coords.GridColumn, rslt.Coords.GridColumn)
			}
			if rslt.Coords.Column != tc.expected.Coords.Column {
				t.Errorf("%s: map_col: expected %d, got %d", tc.name, tc.expected.Coords.Column, rslt.Coords.Column)
			}
			if rslt.Coords.Row != tc.expected.Coords.Row {
				t.Errorf("%s: map_row: expected %d, got %d", tc.name, tc.expected.Coords.Row, rslt.Coords.Row)
			}
			// errors are fragile
			if rslt.Error != nil && tc.expected.Error != nil {
				if errors.Is(tc.err, ast.ErrNoMatch) {
					continue
				} else if !errors.Is(rslt.Error, tc.expected.Error) {
					t.Errorf("%s: error: expected %v, got %v", tc.name, tc.expected.Error, rslt.Error)
				}
			} else if rslt.Error != nil {
				t.Errorf("%s: error: expected nil, got %v", tc.name, rslt.Error)
			} else if tc.expected.Error != nil {
				t.Errorf("%s: error: expected %v, got nil", tc.name, tc.expected.Error)
			}
		}
	}
}

func TestPreviousHexParser(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected ast.PreviousHex_t
		err      error
	}{
		{
			name:  "grid exposed",
			input: `(previous hex = ab 1211)`,
			expected: ast.PreviousHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 2, Column: 12, Row: 11},
			},
		},
		{
			name:  "grid obscured",
			input: `(previous hex = ## 1211)`,
			expected: ast.PreviousHex_t{
				IsObscured: true,
				Coords:     ast.Coords_t{GridRow: 0, GridColumn: 0, Column: 12, Row: 11},
			},
		},
		{
			name:  "n/a",
			input: `(previous hex = n/a)`,
			expected: ast.PreviousHex_t{
				IsNA: true,
			},
		},
		{
			name:     "invalid grid row",
			input:    `(previous hex = 9a 1211)`,
			expected: ast.PreviousHex_t{},
			err:      ast.ErrNoMatch,
		},
		{
			name:     "invalid grid column",
			input:    `(previous hex = a9 1211)`,
			expected: ast.PreviousHex_t{},
			err:      ast.ErrNoMatch,
		},
		{
			name:  "invalid column low",
			input: `(previous hex = aa 0011)`,
			expected: ast.PreviousHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 0, Row: 11},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
		{
			name:  "invalid column high",
			input: `(previous hex = aa 3111)`,
			expected: ast.PreviousHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 31, Row: 11},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
		{
			name:  "invalid row low",
			input: `(previous hex = aa 1100)`,
			expected: ast.PreviousHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 11, Row: 0},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
		{
			name:  "invalid row high",
			input: `(previous hex = aa 1122)`,
			expected: ast.PreviousHex_t{
				Coords: ast.Coords_t{GridRow: 1, GridColumn: 1, Column: 11, Row: 22},
				Error:  ast.ErrInvalidCoordinates,
			},
		},
	} {
		var rslt ast.PreviousHex_t
		var ok bool
		if v, err := units.Parse(tc.name, []byte(tc.input), units.Entrypoint("PreviousHex")); err != nil {
			// errors are fragile
			if tc.err == nil {
				t.Errorf("%s: error: expected nil, got %v", tc.name, err)
			} else if !errors.Is(tc.err, ast.ErrNoMatch) {
				t.Errorf("%s: error: expected %v, got %v", tc.name, tc.err, err)
			}
			continue
		} else if tc.err != nil {
			t.Errorf("%s: error: expected %v, got nil", tc.name, tc.err)
			continue
		} else if rslt, ok = v.(ast.PreviousHex_t); !ok {
			t.Errorf("%s: expected ast.PreviousHex_t, got %T", tc.name, v)
			continue
		}
		if rslt != tc.expected {
			if rslt.IsNA != tc.expected.IsNA {
				t.Errorf("%s: n/a: expected %v, got %v", tc.name, tc.expected.IsNA, rslt.IsNA)
			}
			if rslt.IsObscured != tc.expected.IsObscured {
				t.Errorf("%s: obscured: expected %v, got %v", tc.name, tc.expected.IsObscured, rslt.IsObscured)
			}
			if rslt.Coords.GridRow != tc.expected.Coords.GridRow {
				t.Errorf("%s: grid_row: expected %d, got %d", tc.name, tc.expected.Coords.GridRow, rslt.Coords.GridRow)
			}
			if rslt.Coords.GridColumn != tc.expected.Coords.GridColumn {
				t.Errorf("%s: grid_col: expected %d, got %d", tc.name, tc.expected.Coords.GridColumn, rslt.Coords.GridColumn)
			}
			if rslt.Coords.Column != tc.expected.Coords.Column {
				t.Errorf("%s: map_col: expected %d, got %d", tc.name, tc.expected.Coords.Column, rslt.Coords.Column)
			}
			if rslt.Coords.Row != tc.expected.Coords.Row {
				t.Errorf("%s: map_row: expected %d, got %d", tc.name, tc.expected.Coords.Row, rslt.Coords.Row)
			}
			// errors are fragile
			if rslt.Error != nil && tc.expected.Error != nil {
				if errors.Is(tc.err, ast.ErrNoMatch) {
					continue
				} else if !errors.Is(rslt.Error, tc.expected.Error) {
					t.Errorf("%s: error: expected %v, got %v", tc.name, tc.expected.Error, rslt.Error)
				}
			} else if rslt.Error != nil {
				t.Errorf("%s: error: expected nil, got %v", tc.name, rslt.Error)
			} else if tc.expected.Error != nil {
				t.Errorf("%s: error: expected %v, got nil", tc.name, tc.expected.Error)
			}
		}
	}
}
