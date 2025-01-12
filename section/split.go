// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package section

import (
	"bytes"
	"github.com/playbymail/tribal/is"
	"github.com/playbymail/tribal/norm"
	"sort"
	"strconv"
)

// Split splits the input report into sections.
// Each section contains the header and move data for a single unit.
// All other lines are ignored.
//
// We assume the caller has not done any clean up on the input.
//
// Returns a list of sections sorted by clan id, then by unit id, then by line number.
//
// Warnings:
//   - All input is converted to lowercase to make comparisons easier in future stages.
//   - Report sections sometimes are missing the Status line. We can't depend on it to close out a section.
//   - Sections can contain multiple turn lines because of the missing Status line. When that happens,
//     we capture the additional turn lines and hope that someone eventually reports an error.
func Split(input []byte) (sections []*Section) {
	input = norm.NormalizeSpaces(input)
	input = norm.NormalizeCase(input)
	input = norm.LineEndings(input)

	var section *Section
	for no, line := range bytes.Split(input, []byte{'\n'}) {
		//log.Printf("section: %d: %q\n", no, line)
		if is.UnitHeader(line) {
			// add a new section every time we change units
			section = &Section{
				Id:   len(sections) + 1,
				Line: no + 1,
			}
			section.Lines.Unit = bdup(line)
			sections = append(sections, section)
			// hack the line apart and extract the unit id.
			// we know that the line looks like `tribe tribe_id,` because it is a Unit Header.
			unitUnitId, _, _ := bytes.Cut(line, []byte{','})
			_, section.UnitId, _ = bytes.Cut(unitUnitId, []byte{' '})
			section.ClanId, _ = strconv.Atoi(string(section.UnitId[1:4]))
		} else if section == nil {
			//log.Printf("section: %d: ignoring line %q\n", no, line)
			continue
		} else if is.FleetMovement(line) {
			if section.Lines.FleetMoves == nil {
				section.Lines.FleetMoves = norm.FleetMovement(line)
			}
		} else if is.TribeFollows(line) {
			if section.Lines.UnitFollows == nil {
				section.Lines.UnitFollows = bdup(line)
			}
		} else if is.TribeGoesTo(line) {
			if section.Lines.UnitGoesTo == nil {
				section.Lines.UnitGoesTo = bdup(line)
			}
		} else if is.TribeMovement(line) {
			if section.Lines.UnitMoves == nil {
				section.Lines.UnitMoves = norm.TribeMovement(line)
			}
		} else if is.ScoutLine(line) {
			section.Lines.ScoutLines = append(section.Lines.ScoutLines, norm.ScoutMovement(line))
		} else if is.TurnHeader(line) {
			if section.Lines.Turn == nil {
				section.Lines.Turn = bdup(line)
			}
		} else if is.UnitStatus(line) {
			if section.Lines.Status == nil {
				section.Lines.Status = norm.UnitStatus(line)
			}
			// set `section` to nil to avoid capturing lines between sections.
			section = nil
		}
	}

	// sort the sections by unit then line number.
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Less(sections[j])
	})

	return sections
}

// bdup returns a copy of the slice.
func bdup(b []byte) []byte {
	b2 := make([]byte, len(b))
	copy(b2, b)
	return b2
}
