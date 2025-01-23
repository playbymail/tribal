// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package common

import (
	"bytes"
	"github.com/playbymail/tribal/border"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/terrain"
	"log"
	"strings"
)

// ParseTribeMovement parses the tribe movement line (the "marching" results).
//
// Per the spec, the line should look like this:
//
//	"tribe movement:move" (BACKSLASH MarchSuccess)* (BACKSLASH MarchFail)?
//
// It should be straight-forward to parse, but the report generation process
// has a manual step and that introduces errors that we must deal with. So,
// after we confirm that we have a valid line, we split it into a list of
// segments, using the backslash character as the separator. The list of
// segments should be an optional list of successes followed by an optional
// failure.
//
// We parse the segments and return the list of moves.
func ParseTribeMovement(id ast.UnitId_t, start ast.Coordinates_t, input []byte) (list []*ast.March_t, err error) {
	// split into segments on the backslash
	segments := bytes.Split(input, []byte{'\\'})
	// expect "tribe movement:move" as the first segment
	if len(segments) == 0 || !bytes.Equal(segments[0], []byte("tribe movement:move")) {
		return nil, ast.ErrNotTribeMovementLine
	}
	segments = segments[1:]                       // accept the first segment
	from, previousTerrain := start, terrain.Blank // assign the starting location
	for len(segments) != 0 {
		m, ok := acceptMarchSuccess(id, from, segments[0])
		if !ok {
			break
		}
		list = append(list, m)
		segments, from, previousTerrain = segments[1:], m.To, m.Terrain
	}
	for _, e := range list {
		log.Printf("from %q: march %+v\n", start, *e)
	}

	var failed *ast.March_t
	if failed, segments = acceptMarchFailure(id, from, previousTerrain, segments); failed != nil {
		list = append(list, failed)
	}

	if len(segments) != 0 {
		// accept the excess input.
		// this will have to be presented to the user later.
		m := &ast.March_t{
			Id:        id,
			From:      from,
			Direction: direction.None,
			Terrain:   previousTerrain,
			To:        from,
			Errors:    &ast.MarchErrors_t{},
		}
		for _, seg := range segments {
			m.Errors.ExcessInput = append(m.Errors.ExcessInput, string(seg))
		}
		list = append(list, m)
	}

	return list, nil
}

func acceptMarchFailure(id ast.UnitId_t, from ast.Coordinates_t, fromTerrain terrain.Terrain_e, segments [][]byte) (*ast.March_t, [][]byte) {
	if len(segments) == 0 {
		return nil, segments
	}
	segment := segments[0]
	if match := reCantMove.FindSubmatch(segment); match != nil {
		if ter, ok := terrain.LongTerrainNames[string(match[1])]; ok {
			if dir, ok := direction.LowercaseToEnum[string(match[2])]; ok {
				return &ast.March_t{
					Id:        id,
					From:      from,
					Direction: direction.None,
					To:        from,
					Terrain:   fromTerrain,
					Neighbors: []*ast.Neighbor_t{
						{Terrain: ter, Direction: []direction.Direction_e{dir}},
					},
				}, segments[1:]
			}
		}
	} else if match = reCantMoveWagons.FindSubmatch(segment); match != nil {
		if dir, ok := direction.LowercaseToEnum[string(match[1])]; ok {
			return &ast.March_t{
				Id:        id,
				From:      from,
				Direction: direction.None,
				To:        from,
				Terrain:   fromTerrain,
				Neighbors: []*ast.Neighbor_t{
					{Terrain: terrain.UnknownJungleSwamp, Direction: []direction.Direction_e{dir}},
				},
			}, segments[1:]
		}
	} else if match = reNoFord.FindSubmatch(segment); match != nil {
		if bor, ok := border.LowerCaseToEnum[string(match[1])]; ok {
			if dir, ok := direction.LowercaseToEnum[string(match[2])]; ok {
				return &ast.March_t{
					Id:        id,
					From:      from,
					Direction: direction.None,
					To:        from,
					Terrain:   fromTerrain,
					Borders: []*ast.Border_t{
						{Border: bor, Direction: []direction.Direction_e{dir}},
					},
				}, segments[1:]
			}
		}
	} else if match = reNotEnoughMPs.FindSubmatch(segment); match != nil {
		if dir, ok := direction.LowercaseToEnum[string(match[1])]; ok {
			if ter, ok := terrain.LongTerrainNames[string(match[2])]; ok {
				return &ast.March_t{
					Id:        id,
					From:      from,
					Direction: direction.None,
					To:        from,
					Terrain:   fromTerrain,
					Neighbors: []*ast.Neighbor_t{
						{Terrain: ter, Direction: []direction.Direction_e{dir}},
					},
				}, segments[1:]
			}
		}
	}
	return nil, segments
}

// per the spec
//
//	Direction DASH TerrainCode (COMMA Neighbor)* (COMMA Border)* (COMMA Passage)* (COMMA (SpecialHex | VillageName))?
func acceptMarchSuccess(id ast.UnitId_t, from ast.Coordinates_t, input []byte) (*ast.March_t, bool) {
	log.Printf("accept: success: from %q: input %q\n", from, input)
	dir, ter, rest, ok := AcceptDirectionDashTerrain(input)
	if !ok { // did not find direction-terrain
		log.Printf("accept: success: missing dir-ter %q\n", input)
		return nil, false
	}
	m := &ast.March_t{
		Id:        id,
		From:      from,
		Direction: dir,
		To:        from.Move(dir),
		Terrain:   ter,
	}
	if len(rest) != 0 {
		log.Printf("accept: success: rest %q\n", string(rest))
	}
	input = rest

	// remaining fields are optional
	for len(input) != 0 {
		if input[0] == ' ' || input[0] == ',' {
			input = input[1:]
		} else if elem, rest, ok := acceptNeighbor(input); ok {
			m.Neighbors, input = append(m.Neighbors, elem), rest
		} else if elem, rest, ok := acceptBorder(input); ok {
			m.Borders, input = append(m.Borders, elem), rest
		} else if elem, rest, ok := acceptPassage(input); ok {
			m.Passages, input = append(m.Passages, elem), rest
		} else {
			// we either have a special hex or junk input
			//fmt.Printf("march: input %q\n", input)
			name, rest, _ := bytes.Cut(input, []byte{','})
			if name = bytes.TrimSpace(name); len(name) == 0 {
				// this should be investigated
			} else if m.HexName == nil {
				m.HexName = &ast.HexName_t{Name: strings.Title(string(name))}
			} else {
				if m.Errors == nil {
					m.Errors = &ast.MarchErrors_t{}
				}
				m.Errors.ExcessInput = append(m.Errors.ExcessInput, string(name))
			}
			input = rest
		}
	}

	return m, true
}
