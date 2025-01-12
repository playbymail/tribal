// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package status

import (
	"bytes"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/domains"
	"github.com/playbymail/tribal/resource"
	"github.com/playbymail/tribal/terrain"
	"log"
	"regexp"
)

var (
	reUnitId     = regexp.MustCompile(`^\d{4}([cefg][1-9])?\w`)
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
	s.Terrain = terrainType
	input = rest

	// optional settlement name.
	settlementName, rest, ok := acceptSettlementName(input)
	if ok {
		s.SettlementName = settlementName
		input = rest
	}

	// optional resources
	resourceType, rest, ok := acceptResourceName(input)
	if ok {
		s.Resources = resourceType
		input = rest
	}

	// optional neighbors
	_, rest, ok = acceptNeighborTerrain(input)
	if ok {
		input = rest
	}

	// remaining input is encounters
	for len(input) > 0 {
		unitId, rest, _ := bytes.Cut(input, []byte{' '})
		if !reUnitId.Match(unitId) {
			break
		}
		s.Encounters = append(s.Encounters, domains.UnitId_t(unitId))
		input = rest
	}

	log.Printf("status %q remaining\n", input)

	return &s, nil
}

// accept neighboring terrains.
// these are certain terrain types followed by a list of directions.
func acceptNeighborTerrain(input []byte) ([]domains.Neighbor_t, []byte, bool) {
	if len(input) == 0 || input[0] != ',' {
		return nil, input, false
	}
	var word, rest []byte
	rest = input[1:] // consume the comma
	if idx := bytes.Index(rest, []byte{' '}); idx == -1 {
		// no space found, can't be a neighboring terrain list
		return nil, input, false
	} else {
		word, rest = rest[:idx], rest[idx+1:] // consume the space
	}
	log.Printf("acceptNeighborTerrain: %q %q\n", word, rest)
	enum, ok := terrain.BorderCodes[string(word)]
	if !ok {
		// did not find terrain name
		return nil, input, false
	}
	// we found a terrain name. do we have a list of directions?
	var list []direction.Direction_e
	list, rest, ok = acceptDirections(rest)
	if !ok {
		// did not find directions
		return nil, input, false
	}
	var neighbors []domains.Neighbor_t
	for _, elem := range list {
		neighbors = append(neighbors, domains.Neighbor_t{
			Direction: elem,
			Terrain:   enum,
		})
	}
	log.Printf("acceptNeighborTerrain: returning %q\n", rest)
	return neighbors, rest, true
}

var (
	reDirection = regexp.MustCompile(`^(nw|ne|n|sw|se|s)(?:[ ,]|$)`)
)

// accept list of directions.
// per the spec, the directions are separated by spaces and terminated by a comma (or end of input).
// but because of typos, we'll accept termination by anything other than a direction.
// this is a hack, but it should work because we check for settlement names first.
func acceptDirections(input []byte) ([]direction.Direction_e, []byte, bool) {
	log.Printf("acceptDirections: input %q\n", input)
	match := reDirection.FindSubmatch(input)
	if match == nil {
		// did not find direction name
		return nil, input, false
	}
	word, rest := match[1], input[len(match[1]):] // capture direction and advance
	log.Printf("acceptDirections: word %q rest %q\n", word, rest)
	enum, ok := direction.LowercaseToEnum[string(word)]
	if !ok {
		// did not find direction name
		return nil, input, false
	}
	list := []direction.Direction_e{enum}
	for input = rest; len(input) != 0 && (input[0] == ' '); input = rest {
		input = input[1:] // consume the space
		match = reDirection.FindSubmatch(input)
		if match == nil {
			// did not find direction name
			break
		}
		word, rest = match[1], input[len(match[1]):] // capture direction and advance
		enum, ok = direction.LowercaseToEnum[string(word)]
		if !ok {
			// did not find direction name
			break
		}
		list = append(list, enum)
	}
	return list, input, true
}

// accept resource name followed by comma or end of input
func acceptResourceName(input []byte) (resource.Resource_e, []byte, bool) {
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
	log.Printf("acceptResourceName: %q %q\n", word, rest)
	enum, ok := resource.LongResourceNames[string(word)]
	if !ok {
		// did not find resource name
		return resource.None, input, false
	}
	return enum, rest, true
}

// optional settlement name. this test is horrible.
// if the next item isn't a resource or a neighbor or a unit, then it's a settlement name.
func acceptSettlementName(input []byte) (string, []byte, bool) {
	if len(input) == 0 || input[0] != ',' {
		return "", input, false
	} else if _, _, ok := acceptResourceName(input); ok {
		return "", input, false
	} else if _, _, ok = acceptNeighborTerrain(input); ok {
		return "", input, false
	} else if acceptUnitList(input) {
		return "", input, false
	}
	var word, rest []byte
	rest = input[1:] // consume the comma
	if idx := bytes.Index(rest, []byte{','}); idx == -1 {
		// no comma found, use entire input as terrain name
		word, rest = rest, nil
	} else {
		word, rest = rest[:idx], rest[idx:]
	}
	log.Printf("acceptSettlementName: %q %q\n", word, rest)
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
	log.Printf("acceptTerrainName: %q %q\n", word, rest)
	// is the word a terrain name?
	enum, ok := terrain.LongTerrainNames[string(word)]
	if !ok {
		// did not find terrain name
		return 0, input, false
	}
	return enum, rest, true
}

// accept unit list
func acceptUnitList(input []byte) bool {
	return reUnitId.Match(input)
}
