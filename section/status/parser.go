// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package status

import (
	"bytes"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/domains"
	"github.com/playbymail/tribal/passage"
	"github.com/playbymail/tribal/resource"
	"github.com/playbymail/tribal/terrain"
	"log"
	"regexp"
)

var (
	reUnitStatus = regexp.MustCompile(`^\d{4}([cefg][1-9])? status:`)
)

func Parse(path string, input []byte) (*domains.Status_t, error) {
	var s domains.Status_t

	// expect unit id followed by " status:"
	match := reUnitStatus.Find(input)
	log.Printf("\n\n\nmatch %v input %q\n", string(match), string(input))
	if match == nil {
		return nil, domains.ErrInvalidStatusPrefix
	}
	s.Unit = domains.UnitId_t(string(match))
	input = input[len(match):] // consume the match

	// expect terrain type followed by comma or end of input
	terrainType, rest, ok := acceptTerrainName(input)
	if !ok {
		return nil, domains.ErrMissingTerrainType
	}
	s.Tile.Terrain = terrainType
	input = rest

	// optional settlement name.
	if settlementName, rest, ok := acceptSettlementName(input); ok {
		s.Tile.SettlementName = settlementName
		input = rest
	}

	// optional resources
	s.Tile.Resources, input = acceptOptionalResources(input)
	s.Tile.Neighbors, input = acceptOptionalNeighbors(input)
	s.Tile.Passages, input = acceptOptionalPassages(input)

	// optional encounters
	if encounters, rest, ok := acceptEncounters(input); ok {
		for _, encounter := range encounters {
			s.Tile.Encounters = append(s.Tile.Encounters, encounter)
		}
		input = rest
	}

	log.Printf("status: remaining input %q\n", input)
	if len(input) != 0 {
		s.Errors.ExcessInput = string(input)
	}

	return &s, nil
}

var (
	reTerrainCode = regexp.MustCompile(`^,([a-z]+) `)
)

// accept neighboring terrains.
// these are certain terrain types followed by a list of directions.
func acceptNeighborTerrain(input []byte) ([]domains.Neighbor_t, []byte, bool) {
	match := reTerrainCode.FindSubmatch(input)
	if match == nil { // did not find terrain code
		return nil, input, false
	}
	code, rest := match[1], input[len(match[1])+1:] // capture the terrain code and advance to the delimiter
	enum, ok := terrain.BorderCodes[string(code)]
	if !ok { // did not find terrain code
		return nil, input, false
	}
	var list []direction.Direction_e
	list, rest, ok = acceptDirections(rest)
	if !ok { // did not find terrain code followed by list of directions
		return nil, input, false
	}
	var neighbors []domains.Neighbor_t
	for _, elem := range list {
		neighbors = append(neighbors, domains.Neighbor_t{Terrain: enum, Direction: elem})
	}
	return neighbors, rest, true
}

func acceptOptionalNeighbors(input []byte) (list []domains.Neighbor_t, rest []byte) {
	neighbors, rest, ok := acceptNeighborTerrain(input)
	for ok {
		for _, neighbor := range neighbors {
			list = append(list, neighbor)
		}
		neighbors, rest, ok = acceptNeighborTerrain(rest)
	}
	return list, rest
}

var (
	rePassage = regexp.MustCompile(`^,(canal|ford|pass|river|stone road) `)
)

func acceptPassage(input []byte) ([]domains.Passage_t, []byte, bool) {
	match := rePassage.FindSubmatch(input)
	if match == nil { // did not find passage
		return nil, input, false
	}
	code, rest := match[1], input[len(match[1])+1:] // capture passage and advance to delimiter
	enum, ok := passage.LowerCaseToEnum[string(code)]
	if !ok { // should never happen
		return nil, input, false
	}
	var list []direction.Direction_e
	list, rest, ok = acceptDirections(rest)
	if !ok { // did not find passage followed by directions
		return nil, input, false
	}
	var passages []domains.Passage_t
	for _, elem := range list {
		passages = append(passages, domains.Passage_t{Passage: enum, Direction: elem})
	}
	return passages, rest, true
}

func acceptOptionalPassages(input []byte) (list []domains.Passage_t, rest []byte) {
	passages, rest, ok := acceptPassage(input)
	for ok {
		for _, elem := range passages {
			list = append(list, elem)
		}
		passages, rest, ok = acceptPassage(rest)
	}
	return list, rest
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

// optional settlement name. this test is horrible.
// if the next item isn't a resource or a neighbor or a unit, then it's a settlement name.
func acceptSettlementName(input []byte) (string, []byte, bool) {
	//log.Printf("acceptSettlementName: input %q\n", input)
	if len(input) == 0 || input[0] != ',' {
		return "", input, false
	} else if _, _, ok := acceptResource(input); ok {
		//log.Printf("acceptSettlementName: input %q *** resource\n", input)
		return "", input, false
	} else if _, _, ok = acceptNeighborTerrain(input); ok {
		//log.Printf("acceptSettlementName: input %q *** neighbor\n", input)
		return "", input, false
	} else if _, _, ok = acceptPassage(input); ok {
		//log.Printf("acceptSettlementName: input %q *** passages\n", input)
		return "", input, false
	} else if _, _, ok = acceptEncounters(input); ok {
		//log.Printf("acceptSettlementName: input %q *** encounters\n", input)
		return "", input, false
	}
	//log.Printf("acceptSettlementName: input %q <<--\n", input)
	var word, rest []byte
	rest = input[1:] // consume the comma
	if idx := bytes.Index(rest, []byte{','}); idx == -1 {
		// no comma found, use entire input as settlement name
		word, rest = rest, nil
	} else {
		word, rest = rest[:idx], rest[idx:]
	}
	//log.Printf("acceptSettlementName: word %q rest %q\n", word, rest)
	return string(word), rest, true
}

// accept terrain type followed by comma or end of input
func acceptTerrainName(input []byte) (terrain.Terrain_e, []byte, bool) {
	var word, rest []byte
	idx := bytes.Index(input, []byte{','})
	if idx == -1 {
		// no comma found, use entire input as terrain name
		word, rest = input, nil
	} else {
		word, rest = input[:idx], input[idx:]
	}
	//log.Printf("acceptTerrainName: %q %q\n", word, rest)
	// is the word a terrain name?
	enum, ok := terrain.LongTerrainNames[string(word)]
	if !ok {
		// did not find terrain name
		return 0, input, false
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

var (
	reUnitIdElement = regexp.MustCompile(`^[, ](\d{4}(?:[cefg][1-9])?)(?:[ ,]|$)`)
)

// accept list of unit ids.
// per the spec, the list is (space unit id)* and terminated by a comma (or end of input).
// but because of typos, we'll accept commas or spaces for delimiters and termination by
// anything other than a unit id.
func acceptEncounters(input []byte) ([]domains.UnitId_t, []byte, bool) {
	var list []domains.UnitId_t
	for match := reUnitIdElement.FindSubmatch(input); match != nil; match = reUnitIdElement.FindSubmatch(input) {
		unit, rest := match[1], input[len(match[1])+1:] // capture direction and advance to the delimiter
		list, input = append(list, domains.UnitId_t(unit)), rest
	}
	return list, input, len(list) > 0
}
