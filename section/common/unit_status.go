// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package common

import (
	"github.com/playbymail/tribal/parser/ast"
	"regexp"
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
func ParseUnitStatus(curr ast.Coordinates_t, input []byte) (*ast.Status_t, error) {
	var s ast.Status_t

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
	s.Tile.HexName, input, _ = acceptSpecialHexStatus(input)
	s.Tile.Resources, input = acceptResourceNameList(input)
	s.Tile.Neighbors, input = acceptNeighborList(input)
	s.Tile.Borders, input = acceptBorderList(input)
	s.Tile.Passages, input = acceptPassageList(input)
	s.Tile.Encounters, input = acceptEncounterList(input)

	// if we have something left over, we had invalid input.
	// this will eventually be reported to the user.
	if len(input) != 0 {
		s.Errors.ExcessInput = string(input)
	}

	return &s, nil
}
