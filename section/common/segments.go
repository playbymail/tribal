// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package common implements parser for common segments of movement results.
package common

import (
	"bytes"
	"github.com/playbymail/tribal/border"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/passage"
	"github.com/playbymail/tribal/resource"
	"github.com/playbymail/tribal/terrain"
	"log"
	"regexp"
	"strings"
)

var (
	reBorder               = regexp.MustCompile(`^(canal|river) `)
	reDirectionDashTerrain = regexp.MustCompile(`^([ns][ew]?)-([a-z]{1,3})([ ,]|$)`)
	reDirectionElement     = regexp.MustCompile(`^ ([ns][ew]?)(?: |,|$)`)
	rePassage              = regexp.MustCompile(`^(ford|pass|stone road) `)
	reResource             = regexp.MustCompile(`^find ([a-z]+(?: [a-z]+){0,2})(?: |,|$)`)
	reResourceName         = regexp.MustCompile(`^([a-z]+(?: [a-z]+){0,2})(?: |,|$)`)
	reTerrainCode          = regexp.MustCompile(`^([a-z]{1,3}) `)
	reUnitIdElement        = regexp.MustCompile(`^(\d{4}(?:[cefg][1-9])?)(?:[ ,]|$)`)
)

// acceptBorder returns true if the input starts with a valid border code, and is followed with a list of directions.
//
//	BorderCode (SPACE Direction)+
func acceptBorder(input []byte) (*ast.Border_t, []byte, bool) {
	//log.Printf("accept: border %q\n", input)
	match := reBorder.FindSubmatch(input)
	if match == nil { // did not find border
		return nil, input, false
	}
	code, rest := match[1], input[len(match[1]):] // capture border and advance to delimiter
	enum, ok := border.LowerCaseToEnum[string(code)]
	if !ok { // should never happen
		return nil, input, false
	}
	borders := ast.Border_t{Border: enum}
	//log.Printf("accept: border directions %q\n", rest)
	borders.Direction, rest, ok = AcceptDirectionList(rest)
	if !ok { // did not find passage followed by direction list
		return nil, input, false
	}
	return &borders, rest, true
}

func acceptBorderList(input []byte) ([]*ast.Border_t, []byte) {
	var list []*ast.Border_t
	for elem, rest, ok := acceptBorder(input); ok; elem, rest, ok = acceptBorder(input) {
		list = append(list, elem)
		input = rest
	}
	return list, input
}

//// AcceptBorder returns true if the input starts with a comma,
//// then a valid border code, and is followed with a list of directions.
////
////	COMMA BorderCode (SPACE Direction)+
//func AcceptBorder(input []byte) (*ast.Border_t, []byte, bool) {
//	log.Printf("accept: border %q\n", input)
//	match := reBorder.FindSubmatch(input)
//	if match == nil { // did not find border
//		return nil, input, false
//	}
//	code, rest := match[1], input[len(match[0])-1:] // capture border and advance to delimiter
//	enum, ok := border.LowerCaseToEnum[string(code)]
//	if !ok { // should never happen
//		return nil, input, false
//	}
//	borders := ast.Border_t{Border: enum}
//	log.Printf("accept: border directions %q\n", rest)
//	borders.Direction, rest, ok = AcceptDirectionList(rest)
//	if !ok { // did not find passage followed by direction list
//		return nil, input, false
//	}
//	return &borders, rest, true
//}
//
//func AcceptBorders(input []byte) ([]*ast.Border_t, []byte) {
//	var list []*ast.Border_t
//	for elem, rest, ok := AcceptBorder(input); ok; elem, rest, ok = AcceptBorder(input) {
//		list = append(list, elem)
//		input = rest
//	}
//	return list, input
//}

// AcceptDirectionDashTerrain returns true if the input starts
// with a direction code, a dash, a terrain code and a valid
// separator. If it does, we return the enums for the direction
// and terrain as well.
func AcceptDirectionDashTerrain(input []byte) (dir direction.Direction_e, ter terrain.Terrain_e, rest []byte, ok bool) {
	match := reDirectionDashTerrain.FindSubmatch(input)
	if match == nil { // did not find terrain-direction
		return dir, ter, input, false
	} else if dir, ok = direction.LowercaseToEnum[string(match[1])]; !ok {
		// did not find direction
		return dir, ter, input, false
	} else if ter, ok = terrain.TerrainCodes[string(match[2])]; !ok {
		// did not find terrain code
		return dir, ter, input, false
	}
	sep := match[3]                       // will be empty or a single character
	rest = input[len(match[0])-len(sep):] // advance to delimiter
	return dir, ter, rest, true
}

// AcceptDirectionList returns true if the input starts with a list of directions.
// per the spec, the list is (space direction)* and terminated by a comma (or end of input).
// but because of typos, we'll accept termination by anything other than a direction.
func AcceptDirectionList(input []byte) ([]direction.Direction_e, []byte, bool) {
	var list []direction.Direction_e
	for match := reDirectionElement.FindSubmatch(input); match != nil; match = reDirectionElement.FindSubmatch(input) {
		code, rest := match[1], input[1+len(match[1]):] // capture direction and advance to delimiter
		//log.Printf("acdl: code %q rest %q input %q\n", code, rest, input)
		enum, ok := direction.LowercaseToEnum[string(code)]
		if !ok { // this should never happen
			break
		}
		list, input = append(list, enum), rest
	}
	return list, input, len(list) != 0
}

func acceptEncounter(input []byte) (ast.UnitId_t, []byte, bool) {
	match := reUnitIdElement.FindSubmatch(input)
	if match == nil { // did not find unit id
		return ast.UnitId_t(""), input, false
	}
	unit, rest := match[1], input[len(match[1]):] // capture unit id and advance to the delimiter
	return ast.UnitId_t(unit), rest, true
}

// accept list of unit ids.
// per the spec, the list is (space unit id)* and terminated by a comma (or end of input).
// but because of typos, we'll accept commas or spaces for delimiters and termination by
// anything other than a unit id.
func acceptEncounterList(input []byte) ([]ast.UnitId_t, []byte) {
	var list []ast.UnitId_t
	for match := reUnitIdElement.FindSubmatch(input); match != nil; match = reUnitIdElement.FindSubmatch(input) {
		unit, rest := match[1], input[len(match[1])+1:] // capture unit id and advance to the delimiter
		list, input = append(list, ast.UnitId_t(unit)), rest
	}
	return list, input
}

// acceptHexName reads the input up until the next segment separator and returns it as a hex name.
func acceptHexName(input []byte) (*ast.HexName_t, []byte, bool) {
	if len(input) == 0 {
		return nil, input, false
	}
	var name, rest []byte
	if idx := bytes.Index(input, []byte{','}); idx == -1 {
		name, rest = input, nil // no comma found, use entire input as hex name
	} else {
		name, rest = input[:idx], input[idx:]
	}
	return &ast.HexName_t{Name: strings.Title(string(name))}, rest, true
}

// AcceptNeighbor returns true if the input starts with a valid terrain code, and is followed with a list of directions.
//
//	TerrainCode (SPACE Direction)+
func acceptNeighbor(input []byte) (*ast.Neighbor_t, []byte, bool) {
	match := reTerrainCode.FindSubmatch(input)
	if match == nil { // did not find terrain code
		return nil, input, false
	}
	code, rest := match[1], input[len(match[1]):] // capture the terrain code and advance to the delimiter
	enum, ok := terrain.NeighborCodes[string(code)]
	if !ok { // did not find terrain code
		return nil, input, false
	}
	neighbors := ast.Neighbor_t{Terrain: enum}
	neighbors.Direction, rest, ok = AcceptDirectionList(rest)
	if !ok { // did not find terrain code followed by list of directions
		return nil, input, false
	}
	return &neighbors, rest, true
}

//func AcceptNeighbors(input []byte) ([]*ast.Neighbor_t, []byte) {
//	var list []*ast.Neighbor_t
//	for neighbors, rest, ok := acceptNeighbor(input); ok; neighbors, rest, ok = acceptNeighbor(input) {
//		list = append(list, neighbors)
//		input = rest
//	}
//	return list, input
//}

func acceptNeighborList(input []byte) ([]*ast.Neighbor_t, []byte) {
	var list []*ast.Neighbor_t
	for elem, rest, ok := acceptNeighbor(input); ok; elem, rest, ok = acceptNeighbor(input) {
		list = append(list, elem)
		input = rest
	}
	return list, input
}

// AcceptPassage returns true if the input starts with a passage name and a list of directions.
//
//	Passage (SPACE Direction)+
func acceptPassage(input []byte) (*ast.Passage_t, []byte, bool) {
	match := rePassage.FindSubmatch(input)
	if match == nil { // did not find passage
		return nil, input, false
	}
	code, rest := match[1], input[len(match[1]):] // capture passage and advance to delimiter
	enum, ok := passage.LowerCaseToEnum[string(code)]
	if !ok { // should never happen
		return nil, input, false
	}
	passages := ast.Passage_t{Passage: enum}
	passages.Direction, rest, ok = AcceptDirectionList(rest)
	if !ok { // did not find passage followed by direction list
		return nil, input, false
	}
	return &passages, rest, true
}

func acceptPassageList(input []byte) ([]*ast.Passage_t, []byte) {
	var list []*ast.Passage_t
	for elem, rest, ok := acceptPassage(input); ok; elem, rest, ok = acceptPassage(input) {
		list = append(list, elem)
		input = rest
	}
	return list, input
}

func acceptResource(input []byte) (resource.Resource_e, []byte, bool) {
	match := reResource.FindSubmatch(input)
	if match == nil {
		return resource.None, input, false
	}
	word := match[1]
	//log.Printf("accept: resource %q input %q\n", word, input)
	enum, ok := resource.LongResourceNames[string(word)]
	if !ok { // did not find resource name
		return resource.None, input, false
	}
	input = input[5+len(word):] // consume and advance to the delimiter
	//log.Printf("acceptResourceName: %q %q\n", word, input)
	return enum, input, true
}

func acceptResourceList(input []byte) ([]resource.Resource_e, []byte) {
	var list []resource.Resource_e
	for resourceType, rest, ok := acceptResource(input); ok; resourceType, rest, ok = acceptResource(input) {
		list = append(list, resourceType)
		input = rest
	}
	return list, input
}

func acceptResourceName(input []byte) (resource.Resource_e, []byte, bool) {
	match := reResourceName.FindSubmatch(input)
	if match == nil {
		return resource.None, input, false
	}
	word := match[1]
	//log.Printf("accept: resource %q input %q\n", word, input)
	enum, ok := resource.LongResourceNames[string(word)]
	if !ok { // did not find resource name
		return resource.None, input, false
	}
	input = input[len(word):] // consume and advance to the delimiter
	//log.Printf("acceptResourceName: %q %q\n", word, input)
	return enum, input, true
}

func acceptResourceNameList(input []byte) ([]resource.Resource_e, []byte) {
	var list []resource.Resource_e
	for resourceType, rest, ok := acceptResourceName(input); ok; resourceType, rest, ok = acceptResourceName(input) {
		list = append(list, resourceType)
		input = rest
	}
	return list, input
}

// acceptSpecialHexMarch returns true if the input is a special hex or a village name.
// We all know that this is a lie. It returns true if the input isn't accepted as
// anything else that is allowed in the results from a march (a resource, neighbor,
// border, or passage).
func acceptSpecialHexMarch(input []byte) (*ast.HexName_t, []byte, bool) {
	if len(input) == 0 || input[0] != ',' {
		return nil, input, false
	} else if _, _, ok := acceptResource(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptNeighbor(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptBorder(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptPassage(input); ok {
		return nil, input, false
	}
	input = input[1:] // consume the comma
	return acceptHexName(input)
}

// acceptSpecialHexStatus returns true if the input is a special hex or a village name.
// We all know that this is a lie. It returns true if the input isn't accepted as
// anything else that is allowed in the unit's status line (a resource, neighbor,
// border, passage, or encounter).
func acceptSpecialHexStatus(input []byte) (*ast.HexName_t, []byte, bool) {
	log.Printf("acceptSpecialHex %q\n", input)
	if len(input) == 0 || input[0] != ',' {
		return nil, input, false
	} else if _, _, ok := acceptResourceName(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptNeighbor(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptBorder(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptPassage(input); ok {
		return nil, input, false
	} else if _, _, ok = acceptEncounter(input); ok {
		return nil, input, false
	}
	input = input[1:] // consume the comma
	return acceptHexName(input)
}

// acceptTerrainName followed by comma or end of input
func acceptTerrainName(input []byte) (terrain.Terrain_e, []byte, bool) {
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
