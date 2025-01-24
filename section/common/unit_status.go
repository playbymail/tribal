// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package common

import (
	"bytes"
	"github.com/playbymail/tribal/parser/ast"
	"regexp"
	"strings"
)

var (
	reStatusPrefix = regexp.MustCompile(`^(\d{4}([cefg][1-9])?) status:`)
)

// ParseUnitStatus parses the status of a unit.
// The status line is required, but it's sometimes missing in setup reports.
//
// Per the spec, the line should look like this:
//
//	UnitId "status:" TerrainName (COMMA (SpecialHex | VillageName))? (COMMA Resources)* (COMMA Neighbor)* (COMMA Border)* (COMMA Passages)* (COMMA Units)*
func ParseUnitStatus(turn *ast.Turn_t, curr ast.Coordinates_t, input []byte) (*ast.Status_t, error) {
	s := ast.Status_t{
		Turn: turn,
	}

	// expect unit id followed by " status:"
	if match := reStatusPrefix.FindSubmatch(input); match == nil {
		return nil, ast.ErrNotUnitStatusLine
	} else {
		s.Unit = ast.UnitId_t(match[1])
		input = input[len(match[0]):] // consume the match
	}

	// expect terrain name followed by comma or end of input
	if terrainType, rest, ok := acceptTerrainName(input); !ok {
		return nil, ast.ErrMissingTerrainType
	} else {
		s.Tile.Terrain = terrainType
		input = rest
	}
	// tile will inherit this unit's current location
	s.Tile.Coordinates = curr

	// remaining fields are optional
	for len(input) != 0 {
		if input[0] == ' ' || input[0] == ',' {
			input = input[1:]
		} else if elem, rest, ok := acceptResourceName(input); ok {
			s.Tile.Resources, input = append(s.Tile.Resources, elem), rest
		} else if elem, rest, ok := acceptNeighbor(input); ok {
			s.Tile.Neighbors, input = append(s.Tile.Neighbors, elem), rest
		} else if elem, rest, ok := acceptBorder(input); ok {
			s.Tile.Borders, input = append(s.Tile.Borders, elem), rest
		} else if elem, rest, ok := acceptPassage(input); ok {
			s.Tile.Passages, input = append(s.Tile.Passages, elem), rest
		} else if elem, rest, ok := acceptEncounter(input); ok {
			s.Tile.Encounters, input = append(s.Tile.Encounters, elem), rest
		} else {
			// we either have a special hex or junk input
			//fmt.Printf("status: input %q\n", input)
			name, rest, _ := bytes.Cut(input, []byte{','})
			if name = bytes.TrimSpace(name); len(name) == 0 {
				// this should be investigated
			} else if s.Tile.HexName == nil {
				s.Tile.HexName = &ast.HexName_t{Name: strings.Title(string(name))}
			} else {
				if s.Errors == nil {
					s.Errors = &ast.StatusErrors_t{}
				}
				s.Errors.ExcessInput = append(s.Errors.ExcessInput, string(name))
			}
			input = rest
		}
		//s.Tile.Resources, input = acceptResourceNameList(input)
		//s.Tile.Neighbors, input = acceptNeighborList(input)
		//s.Tile.Borders, input = acceptBorderList(input)
		//s.Tile.Passages, input = acceptPassageList(input)
		//s.Tile.Encounters, input = acceptEncounterList(input)
	}

	return &s, nil
}
