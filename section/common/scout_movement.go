// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package common

import (
	"bytes"
	"github.com/playbymail/tribal/border"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/item"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/terrain"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	reScoutPatrol = regexp.MustCompile("^scout ([1-8]):scout(?: |,|$)")
)

// ParseScoutMovement parses the scout patrol line (the "patrol" results).
//
// Per the spec, the line should look like this:
//
//	"scout" ScoutId ":scout" (BACKSLASH PatrolSuccess)* (BACKSLASH PatrolFail)?
//
// It should be straight-forward to parse, but the report generation process
// has a manual step and that introduces errors that we must deal with. So,
// after we confirm that we have a valid line, we split it into a list of
// segments, using the backslash character as the separator. The list of
// segments should be an optional list of successes followed by an optional
// failure.
//
// We parse the segments and return the list of patrol results.
func ParseScoutMovement(turn *ast.Turn_t, id ast.UnitId_t, start ast.Coordinates_t, input []byte) (list []*ast.Patrol_t, err error) {
	// split into segments on the backslash
	segments := bytes.Split(input, []byte{'\\'})
	// expect "scout" ScoutId ":scout" as the first segment
	if len(segments) == 0 || reScoutPatrol.FindSubmatch(segments[0]) == nil {
		log.Printf("psm: turn %d unit %q input %q\n", turn, id, input)
		return nil, ast.ErrNotScoutPatrolLine
	}
	match := reScoutPatrol.FindSubmatch(segments[0])
	if match == nil {
		log.Printf("psm: turn %d unit %q input %q\n", turn, id, input)
		return nil, ast.ErrNotScoutPatrolLine
	}
	patrolId := int(match[1][0] - '0')
	//log.Printf("scout: %d: from %q: input %q\n", patrolId, start, input)

	segments = segments[1:]                       // accept the first segment
	from, previousTerrain := start, terrain.Blank // assign the starting location
	_ = previousTerrain

	// big loop should process all the things, unfortunately
	//if turn == 19 && id == "0163" && patrolId == 1 {
	//	fmt.Printf("sp input %q\n", input)
	//}
	for _, seg := range segments {
		//if turn == 19 && id == "0163" && patrolId == 1 {
		//	fmt.Printf("sp seg %q\n", seg)
		//}
		if ps, ok := acceptPatrolSuccess(turn, id, patrolId, from, seg); ok {
			list = append(list, ps)
		} else if ps, ok := acceptPatrolFailure(turn, id, patrolId, from, previousTerrain, seg); ok {
			list = append(list, ps)
		} else if ps, ok := acceptPatrolFound(turn, id, patrolId, from, previousTerrain, seg); ok {
			list = append(list, ps)
		} else {
			// if we get to here, we've got a segment that we don't know how to process
			list = append(list, &ast.Patrol_t{
				Turn:   turn,
				Id:     id,
				Patrol: patrolId,
				From:   from,
				To:     from,
				Errors: &ast.PatrolErrors_t{ExcessInput: []string{string(seg)}},
			})
		}
	}

	return list, nil
}

var (
	reCantMove       = regexp.MustCompile(`^can't move on ([a-z]+(?: [a-z]+){0,2}) to ([ns][ew]?) of hex$`)
	reCantMoveWagons = regexp.MustCompile(`^cannot move wagons into swamp/jungle hill to ([ns][ew]?) of hex$`)
	reFindQtyItem    = regexp.MustCompile(`^find (\d+) ([a-z]+(?: [a-z]+){0,2})$`)
	reNoFord         = regexp.MustCompile(`^no ford on (canal|river) to ([ns][ew]?) of hex$`)
	reNotEnoughMPs   = regexp.MustCompile(`^not enough m\.p's to move to ([ns][ew]?) into ([a-z]+(?: [a-z]+){0,2})$`)
)

func acceptPatrolFailure(turn *ast.Turn_t, id ast.UnitId_t, patrolId int, from ast.Coordinates_t, fromTerrain terrain.Terrain_e, input []byte) (*ast.Patrol_t, bool) {
	if match := reCantMove.FindSubmatch(input); match != nil {
		if ter, ok := terrain.LongTerrainNames[string(match[1])]; ok {
			if dir, ok := direction.LowercaseToEnum[string(match[2])]; ok {
				return &ast.Patrol_t{
					Turn:      turn,
					Id:        id,
					Patrol:    patrolId,
					From:      from,
					Direction: direction.None,
					To:        from,
					Terrain:   fromTerrain,
					Neighbors: []*ast.Neighbor_t{
						{Terrain: ter, Direction: []direction.Direction_e{dir}},
					},
				}, true
			}
		}
	} else if match = reNoFord.FindSubmatch(input); match != nil {
		if bor, ok := border.LowerCaseToEnum[string(match[1])]; ok {
			if dir, ok := direction.LowercaseToEnum[string(match[2])]; ok {
				return &ast.Patrol_t{
					Turn:      turn,
					Id:        id,
					Patrol:    patrolId,
					From:      from,
					Direction: direction.None,
					To:        from,
					Terrain:   fromTerrain,
					Borders: []*ast.Border_t{
						{Border: bor, Direction: []direction.Direction_e{dir}},
					},
				}, true
			}
		}
	} else if match = reNotEnoughMPs.FindSubmatch(input); match != nil {
		if dir, ok := direction.LowercaseToEnum[string(match[1])]; ok {
			if ter, ok := terrain.LongTerrainNames[string(match[2])]; ok {
				return &ast.Patrol_t{
					Turn:      turn,
					Id:        id,
					Patrol:    patrolId,
					From:      from,
					Direction: direction.None,
					To:        from,
					Terrain:   fromTerrain,
					Neighbors: []*ast.Neighbor_t{
						{Terrain: ter, Direction: []direction.Direction_e{dir}},
					},
				}, true
			}
		}
	}
	return nil, false
}

func acceptPatrolFound(turn *ast.Turn_t, id ast.UnitId_t, patrolId int, from ast.Coordinates_t, ter terrain.Terrain_e, input []byte) (*ast.Patrol_t, bool) {
	if bytes.HasPrefix(input, []byte(`no groups located`)) {
		input = input[17:] // consume prefix
		ps := &ast.Patrol_t{
			Turn:      turn,
			Id:        id,
			Patrol:    patrolId,
			From:      from,
			Direction: direction.None,
			To:        from,
			Terrain:   ter,
		}
		// if we have something left over, we had invalid input.
		// this will eventually be reported to the user.
		if len(input) != 0 {
			//log.Printf("accept: %q: scout %d: patrol excess %q\n", id, patrolId, string(input))
			if ps.Errors == nil {
				ps.Errors = &ast.PatrolErrors_t{}
			}
			ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(input))
		}
		return ps, true
	} else if bytes.HasPrefix(input, []byte(`nothing of interest found`)) {
		input = input[25:] // consume prefix
		ps := &ast.Patrol_t{
			Turn:      turn,
			Id:        id,
			Patrol:    patrolId,
			From:      from,
			Direction: direction.None,
			To:        from,
			Terrain:   ter,
		}
		// if we have something left over, we had invalid input.
		// this will eventually be reported to the user.
		if len(input) != 0 {
			//log.Printf("accept: %q: scout %d: patrol excess %q\n", id, patrolId, string(input))
			if ps.Errors == nil {
				ps.Errors = &ast.PatrolErrors_t{}
			}
			ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(input))
		}
		return ps, true
	} else if bytes.HasPrefix(input, []byte(`patrolled and found `)) {
		input = input[20:] // consume prefix up to and including the delimiting space
		ps := &ast.Patrol_t{
			Turn:      turn,
			Id:        id,
			Patrol:    patrolId,
			From:      from,
			Direction: direction.None,
			To:        from,
			Terrain:   ter,
		}
		ps.Encounters, input = acceptEncounterList(input)
		// if we have something left over, we had invalid input.
		// this will eventually be reported to the user.
		if len(input) != 0 {
			//log.Printf("accept: %q: scout %d: patrol excess %q\n", id, patrolId, string(input))
			if ps.Errors == nil {
				ps.Errors = &ast.PatrolErrors_t{}
			}
			ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(input))
		}
		return ps, true
	} else if match := reFindQtyItem.FindSubmatch(input); match != nil {
		if enum, ok := item.LowerCaseName[string(match[2])]; ok {
			if qty, err := strconv.Atoi(string(match[1])); err == nil && qty > 0 {
				return &ast.Patrol_t{
					Turn:      turn,
					Id:        id,
					Patrol:    patrolId,
					From:      from,
					Direction: direction.None,
					To:        from,
					Terrain:   ter,
					Items:     []ast.Item_t{{Item: enum, Quantity: qty}},
				}, true
			}
		}
	}
	return nil, false
}

func acceptPatrolSuccess(turn *ast.Turn_t, id ast.UnitId_t, patrolId int, from ast.Coordinates_t, input []byte) (*ast.Patrol_t, bool) {
	//log.Printf("accept: success: from %q: input %q\n", from, input)
	dir, ter, rest, ok := AcceptDirectionDashTerrain(input)
	if !ok { // did not find direction-terrain
		//log.Printf("accept: success: missing dir-ter %q\n", input)
		return nil, false
	}
	ps := &ast.Patrol_t{
		Turn:      turn,
		Id:        id,
		Patrol:    patrolId,
		From:      from,
		Direction: dir,
		To:        from.Move(dir),
		Terrain:   ter,
	}
	input = rest

	// remaining fields are optional
	for len(input) != 0 {
		if input[0] == ' ' || input[0] == ',' {
			input = input[1:]
		} else if elem, rest, ok := acceptNeighbor(input); ok {
			ps.Neighbors, input = append(ps.Neighbors, elem), rest
		} else if elem, rest, ok := acceptBorder(input); ok {
			ps.Borders, input = append(ps.Borders, elem), rest
		} else if elem, rest, ok := acceptPassage(input); ok {
			ps.Passages, input = append(ps.Passages, elem), rest
		} else if elem, rest, ok := acceptResource(input); ok {
			ps.Resources, input = append(ps.Resources, elem), rest
		} else if elem, rest, ok := acceptEncounter(input); ok {
			ps.Encounters, input = append(ps.Encounters, elem), rest
		} else {
			// we either have a special hex or junk input
			//fmt.Printf("patrol: input %q\n", input)
			name, rest, _ := bytes.Cut(input, []byte{','})
			if name = bytes.TrimSpace(name); len(name) == 0 {
				// this should be investigated
			} else if ps.HexName == nil {
				ps.HexName = &ast.HexName_t{Name: strings.Title(string(name))}
			} else {
				if ps.Errors == nil {
					ps.Errors = &ast.PatrolErrors_t{}
				}
				ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(name))
			}
			input = rest
		}
	}

	//// remaining fields are optional
	//ps.Neighbors, input = acceptNeighborList(input)
	//ps.Borders, input = acceptBorderList(input)
	//ps.Passages, input = acceptPassageList(input)
	//ps.Resources, input = acceptResourceList(input)
	//// special hex or village name before encounters? ick?
	//ps.Encounters, input = acceptEncounterList(input)
	//
	//// if we have something left over, we had invalid input.
	//// this will eventually be reported to the user.
	//if len(input) != 0 {
	//	log.Printf("accept: %q: scout %d: success excess %q\n", unit, id, string(input))
	//	if ps.Errors == nil {
	//		ps.Errors = &ast.PatrolErrors_t{}
	//	}
	//	ps.Errors.ExcessInput = string(input)
	//}

	return ps, true
}
