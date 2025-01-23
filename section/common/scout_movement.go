// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package common

import (
	"bytes"
	"github.com/playbymail/tribal/border"
	"github.com/playbymail/tribal/direction"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/terrain"
	"regexp"
	"strings"
)

var (
	reScoutPatrol = regexp.MustCompile("^scout ([1-8]):scout(?: |$)")
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
func ParseScoutMovement(id ast.UnitId_t, start ast.Coordinates_t, input []byte) (list []*ast.Patrol_t, err error) {
	// split into segments on the backslash
	segments := bytes.Split(input, []byte{'\\'})
	// expect "scout" ScoutId ":scout" as the first segment
	if len(segments) == 0 || reScoutPatrol.FindSubmatch(segments[0]) == nil {
		return nil, ast.ErrNotScoutPatrolLine
	}
	match := reScoutPatrol.FindSubmatch(segments[0])
	if match == nil {
		return nil, ast.ErrNotScoutPatrolLine
	}
	patrolId := int(match[1][0] - '0')
	//log.Printf("scout: %d: from %q: input %q\n", patrolId, start, input)

	segments = segments[1:]                       // accept the first segment
	from, previousTerrain := start, terrain.Blank // assign the starting location
	_ = previousTerrain

	for len(segments) != 0 {
		ps, ok := acceptPatrolSuccess(id, patrolId, from, segments[0])
		if !ok {
			break
		}
		list = append(list, ps)
		segments, from, previousTerrain = segments[1:], ps.To, ps.Terrain
	}

	var failed *ast.Patrol_t
	if failed, segments = acceptPatrolFailure(id, patrolId, from, previousTerrain, segments); failed != nil {
		list = append(list, failed)
	}

	if len(segments) != 0 { // see if we have a patrol
		if ps, ok := acceptPatrolFound(id, patrolId, from, previousTerrain, segments[0]); ok {
			list = append(list, ps)
			segments = segments[1:]
		}
	}

	if len(segments) != 0 {
		// accept the excess input.
		// this will have to be presented to the user later.
		ps := &ast.Patrol_t{
			Id:     id,
			Patrol: patrolId,
			From:   from,
			To:     from,
			Errors: &ast.PatrolErrors_t{},
		}
		for _, seg := range segments {
			ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(seg))
		}
		list = append(list, ps)
	}

	return list, nil
}

var (
	reCantMove       = regexp.MustCompile(`^can't move on ([a-z]+(?: [a-z]+){0,2}) to ([ns][ew]?) of hex$`)
	reCantMoveWagons = regexp.MustCompile(`^cannot move wagons into swamp/jungle hill to ([ns][ew]?) of hex$`)
	reNoFord         = regexp.MustCompile(`^no ford on (canal|river) to ([ns][ew]?) of hex$`)
	reNotEnoughMPs   = regexp.MustCompile(`^not enough m\.p's to move to ([ns][ew]?) into ([a-z]+(?: [a-z]+){0,2})$`)
)

func acceptPatrolFailure(id ast.UnitId_t, patrolId int, from ast.Coordinates_t, fromTerrain terrain.Terrain_e, segments [][]byte) (*ast.Patrol_t, [][]byte) {
	if len(segments) == 0 {
		return nil, segments
	}
	segment := segments[0]
	if match := reCantMove.FindSubmatch(segment); match != nil {
		if ter, ok := terrain.LongTerrainNames[string(match[1])]; ok {
			if dir, ok := direction.LowercaseToEnum[string(match[2])]; ok {
				return &ast.Patrol_t{
					Id:        id,
					Patrol:    patrolId,
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
	} else if match = reNoFord.FindSubmatch(segment); match != nil {
		if bor, ok := border.LowerCaseToEnum[string(match[1])]; ok {
			if dir, ok := direction.LowercaseToEnum[string(match[2])]; ok {
				return &ast.Patrol_t{
					Id:        id,
					Patrol:    patrolId,
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
				return &ast.Patrol_t{
					Id:        id,
					Patrol:    patrolId,
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

func acceptPatrolFound(unit ast.UnitId_t, id int, from ast.Coordinates_t, ter terrain.Terrain_e, input []byte) (*ast.Patrol_t, bool) {
	if bytes.HasPrefix(input, []byte(`no groups located`)) {
		input = input[17:] // consume prefix
		ps := &ast.Patrol_t{
			Id:        unit,
			Patrol:    id,
			From:      from,
			Direction: direction.None,
			To:        from,
			Terrain:   ter,
		}
		// if we have something left over, we had invalid input.
		// this will eventually be reported to the user.
		if len(input) != 0 {
			//log.Printf("accept: %q: scout %d: patrol excess %q\n", unit, id, string(input))
			if ps.Errors == nil {
				ps.Errors = &ast.PatrolErrors_t{}
			}
			ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(input))
		}
		return ps, true
	} else if bytes.HasPrefix(input, []byte(`nothing of interest found`)) {
		input = input[25:] // consume prefix
		ps := &ast.Patrol_t{
			Id:        unit,
			Patrol:    id,
			From:      from,
			Direction: direction.None,
			To:        from,
			Terrain:   ter,
		}
		// if we have something left over, we had invalid input.
		// this will eventually be reported to the user.
		if len(input) != 0 {
			//log.Printf("accept: %q: scout %d: patrol excess %q\n", unit, id, string(input))
			if ps.Errors == nil {
				ps.Errors = &ast.PatrolErrors_t{}
			}
			ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(input))
		}
		return ps, true
	} else if bytes.HasPrefix(input, []byte(`patrolled and found `)) {
		input = input[19:] // consume prefix up to the delimiting space
		ps := &ast.Patrol_t{
			Id:        unit,
			Patrol:    id,
			From:      from,
			Direction: direction.None,
			To:        from,
			Terrain:   ter,
		}
		ps.Encounters, input = acceptEncounterList(input)
		// if we have something left over, we had invalid input.
		// this will eventually be reported to the user.
		if len(input) != 0 {
			//log.Printf("accept: %q: scout %d: patrol excess %q\n", unit, id, string(input))
			if ps.Errors == nil {
				ps.Errors = &ast.PatrolErrors_t{}
			}
			ps.Errors.ExcessInput = append(ps.Errors.ExcessInput, string(input))
		}
		return ps, true
	}
	return nil, false
}

func acceptPatrolSuccess(unit ast.UnitId_t, id int, from ast.Coordinates_t, input []byte) (*ast.Patrol_t, bool) {
	//log.Printf("accept: success: from %q: input %q\n", from, input)
	dir, ter, rest, ok := AcceptDirectionDashTerrain(input)
	if !ok { // did not find direction-terrain
		//log.Printf("accept: success: missing dir-ter %q\n", input)
		return nil, false
	}
	ps := &ast.Patrol_t{
		Id:        unit,
		Patrol:    id,
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
