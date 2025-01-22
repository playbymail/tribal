// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package status

import (
	"bytes"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/passage"
	"github.com/playbymail/tribal/resource"
	"github.com/playbymail/tribal/terrain"
	"regexp"
	"strings"
)

var (
	reStatusPrefix = regexp.MustCompile(`^(\d{4}([cefg][1-9])?) status:`)
)

// Parse parses the status of a unit.
// The status line is required, but it's sometimes missing in setup reports.
//
// Per the spec, the line should look like this:
//
//	UnitId "status:" TerrainName (COMMA (SpecialHex | VillageName))? (COMMA Resources)* (COMMA Neighbor)* (COMMA Border)* (COMMA Passages)* (COMMA Units)*
func Parse(curr ast.Coordinates_t, input []byte) (*ast.Status_t, error) {
	var s ast.Status_t

	// expect unit id followed by " status:"
	if match := reStatusPrefix.FindSubmatch(input); match == nil {
		return nil, ast.ErrNotUnitStatusLine
	} else {
		s.Unit = ast.UnitId_t(match[1])
		input = input[len(match[0]):] // consume the match
	}

	// expect terrain name followed by comma or end of input
	if terrainType, rest, ok := expectTerrainName(input); !ok {
		return nil, ast.ErrMissingTerrainType
	} else {
		s.Tile.Terrain = terrainType
		input = rest
	}
	// tile will inherit this unit's current location
	s.Tile.Coordinates = curr

	// remaining fields are optional
	s.Tile.HexName, input = acceptOptionalSpecialHexName(input)
	s.Tile.Resources, input = acceptOptionalResources(input)
	s.Tile.Neighbors, input = acceptOptionalNeighbors(input)
	s.Tile.Passages, input = acceptOptionalPassages(input)
	s.Tile.Encounters, input = acceptOptionalEncounters(input)

	// if we have something left over, we had invalid input.
	// this will eventually be reported to the user.
	if len(input) != 0 {
		s.Errors.ExcessInput = string(input)
	}

	return &s, nil
}

var (
	reUnitIdElement = regexp.MustCompile(`^[, ](\d{4}(?:[cefg][1-9])?)(?:[ ,]|$)`)
)

func acceptEncounter(input []byte) (ast.UnitId_t, []byte, bool) {
	match := reUnitIdElement.FindSubmatch(input)
	if match == nil { // did not find unit id
		return "", input, false
	}
	unit, rest := match[1], input[len(match[1])+1:] // capture unit id and advance to the delimiter
	return ast.UnitId_t(unit), rest, true
}

func acceptOptionalEncounters(input []byte) (list []ast.UnitId_t, rest []byte) {
	unit, rest, ok := acceptEncounter(input)
	for ok {
		list = append(list, unit)
		unit, rest, ok = acceptEncounter(rest)
	}
	return list, rest
}

var (
	reTerrainCode = regexp.MustCompile(`^,([a-z]+) `)
)

// accept neighboring terrains.
// these are certain terrain types followed by a list of directions.
func acceptNeighborTerrain(input []byte) (*ast.Neighbor_t, []byte, bool) {
	match := reTerrainCode.FindSubmatch(input)
	if match == nil { // did not find terrain code
		return nil, input, false
	}
	code, rest := match[1], input[len(match[1])+1:] // capture the terrain code and advance to the delimiter
	enum, ok := terrain.NeighborCodes[string(code)]
	if !ok { // did not find terrain code
		return nil, input, false
	}
	neighbors := ast.Neighbor_t{Terrain: enum}
	neighbors.Direction, rest, ok = acceptDirections(rest)
	if !ok { // did not find terrain code followed by list of directions
		return nil, input, false
	}
	return &neighbors, rest, true
}

func acceptOptionalNeighbors(input []byte) ([]*ast.Neighbor_t, []byte) {
	var list []*ast.Neighbor_t
	for neighbors, rest, ok := acceptNeighborTerrain(input); ok; neighbors, rest, ok = acceptNeighborTerrain(input) {
		list = append(list, neighbors)
		input = rest
	}
	return list, input
}

var (
	rePassage = regexp.MustCompile(`^,(canal|ford|pass|river|stone road) `)
)

func acceptPassage(input []byte) (*ast.Passage_t, []byte, bool) {
	match := rePassage.FindSubmatch(input)
	if match == nil { // did not find passage
		return nil, input, false
	}
	code, rest := match[1], input[len(match[1])+1:] // capture passage and advance to delimiter
	enum, ok := passage.LowerCaseToEnum[string(code)]
	if !ok { // should never happen
		return nil, input, false
	}
	passages := ast.Passage_t{Passage: enum}
	passages.Direction, rest, ok = acceptDirections(rest)
	if !ok { // did not find passage followed by direction list
		return nil, input, false
	}
	return &passages, rest, true
}

func acceptOptionalPassages(input []byte) ([]*ast.Passage_t, []byte) {
	var list []*ast.Passage_t
	for passages, rest, ok := acceptPassage(input); ok; passages, rest, ok = acceptPassage(input) {
		list = append(list, passages)
		input = rest
	}
	return list, input
}

// accept resource name followed by comma or end of input
func acceptResource(input []byte) (resource.Resource_e, []byte, bool) {
	if len(input) == 0 || input[0] != ',' {
		return resource.None, input, false
	}
	var word, rest []byte
	rest = input[1:] // consume the comma
	if idx := bytes.Index(rest, []byte{','}); idx == -1 {
		// no comma found, use entire input as resource name
		word, rest = rest, nil
	} else {
		word, rest = rest[:idx], rest[idx:]
	}
	//log.Printf("acceptResourceName: %q %q\n", word, rest)
	enum, ok := resource.LongResourceNames[string(word)]
	if !ok {
		// did not find resource name
		return resource.None, input, false
	}
	return enum, rest, true
}

func acceptOptionalResources(input []byte) (list []resource.Resource_e, rest []byte) {
	resourceType, rest, ok := acceptResource(input)
	for ok {
		list = append(list, resourceType)
		resourceType, rest, ok = acceptResource(rest)
	}
	return list, rest
}

// hex name. this test is horrible.
// if the next item isn't something else, then it's a special hex or village name.
func acceptHexName(input []byte) (*ast.HexName_t, []byte, bool) {
	if len(input) == 0 || input[0] != ',' {
		return nil, input, false
	} else if _, _, ok := acceptResource(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptNeighborTerrain(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptPassage(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptEncounter(input); ok {
		return nil, input, false
	}
	rest := input[1:] // consume the comma
	var name []byte
	if idx := bytes.Index(rest, []byte{','}); idx == -1 {
		name, rest = rest, nil // no comma found, use entire input as settlement name
	} else {
		name, rest = rest[:idx], rest[idx:]
	}
	return &ast.HexName_t{Name: strings.Title(string(name))}, rest, true
}

func acceptOptionalSpecialHexName(input []byte) (*ast.HexName_t, []byte) {
	name, rest, ok := acceptHexName(input)
	if !ok {
		return name, input
	}
	return name, rest
}

// expect terrain name followed by comma or end of input
func expectTerrainName(input []byte) (terrain.Terrain_e, []byte, bool) {
	var name, rest []byte
	if idx := bytes.Index(input, []byte{','}); idx == -1 {
		name, rest = input, nil // no comma found, use entire input as terrain name
	} else {
		name, rest = input[:idx], input[idx:]
	}
	enum, ok := terrain.LongTerrainNames[string(name)]
	if !ok { // did not find terrain name
		return terrain.Blank, input, false
	}
	return enum, rest, true
}

var (
	reDirectionElement = regexp.MustCompile(`^ (nw|ne|n|sw|se|s)(?:[ ,]|$)`)
)

// accept list of directions.
// per the spec, the list is (space direction)* and terminated by a comma (or end of input).
// but because of typos, we'll accept termination by anything other than a direction.
func acceptDirections(input []byte) ([]direction.Direction_e, []byte, bool) {
	var list []direction.Direction_e
	for match := reDirectionElement.FindSubmatch(input); match != nil; match = reDirectionElement.FindSubmatch(input) {
		code, rest := match[1], input[len(match[1])+1:] // capture direction and advance to delimiter
		enum, ok := direction.LowercaseToEnum[string(code)]
		if !ok { // this should never happen
			break
		}
		list, input = append(list, enum), rest
	}
	return list, input, len(list) != 0
}

// accept list of unit ids.
// per the spec, the list is (space unit id)* and terminated by a comma (or end of input).
// but because of typos, we'll accept commas or spaces for delimiters and termination by
// anything other than a unit id.
func acceptEncounters(input []byte) ([]ast.UnitId_t, []byte, bool) {
	var list []ast.UnitId_t
	for match := reUnitIdElement.FindSubmatch(input); match != nil; match = reUnitIdElement.FindSubmatch(input) {
		unit, rest := match[1], input[len(match[1])+1:] // capture unit id and advance to the delimiter
		list, input = append(list, ast.UnitId_t(unit)), rest
	}
	return list, input, len(list) > 0
}
