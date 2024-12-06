// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package section

import (
	"bytes"
	"github.com/playbymail/tribal/is"
	"github.com/playbymail/tribal/norm"
	"regexp"
)

type Section struct {
	Line   int // line number in the original input
	Id     int // section number, starting at 1
	Header []byte
	Turn   []byte
	Moves  struct {
		Movement []byte
		Follows  []byte
		GoesTo   []byte
		Fleet    []byte
		Scouts   [][]byte
	}
	Status []byte
}

// Split splits the input report into sections.
// Each section contains the header and move data for a single unit.
// All other lines are ignored.
//
// All input is converted to lowercase to make comparisons easier in future stages.
func Split(input [][]byte) (sections []*Section) {
	var section *Section
	for no, line := range input {
		// we can't assume that the caller has cleaned up the input.
		//log.Printf("section: %d: %q\n", no, line)
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// force to lowercase and compress whitespace
		line = norm.NormalizeSpaces(bytes.ToLower(line))
		//log.Printf("section: %d: %q\n", no, line)
		if is.UnitHeader(line) {
			section = &Section{
				Id:     len(sections) + 1,
				Line:   no + 1,
				Header: bdup(line),
			}
			sections = append(sections, section)
		} else if section == nil {
			continue
		} else if is.FleetMovement(line) {
			section.Moves.Fleet = line
		} else if is.TribeFollows(line) {
			section.Moves.Follows = line
		} else if is.TribeGoesTo(line) {
			section.Moves.GoesTo = line
		} else if is.TribeMovement(line) {
			section.Moves.Movement = line
		} else if is.ScoutLine(line) {
			section.Moves.Scouts = append(section.Moves.Scouts, line)
		} else if is.TurnHeader(line) {
			section.Turn = line
		} else if is.UnitStatus(line) {
			section.Status = line
		}
	}
	return sections
}

var (
	rxTurnHeader = regexp.MustCompile(`^current turn (\d{3,4})-(\d{1,2})\(#\d+\),`)
)

// bdup returns a copy of the slice.
func bdup(b []byte) []byte {
	b2 := make([]byte, len(b))
	copy(b2, b)
	return b2
}
