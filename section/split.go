// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package section

import (
	"bytes"
	"github.com/playbymail/tribal/is"
	"github.com/playbymail/tribal/norm"
	"sort"
	"strconv"
)

type Section struct {
	Line   int    // line number in the original input
	Id     int    // section number, starting at 1
	ClanId int    // derived from the header
	UnitId []byte // taken from the header
	Header []byte // text of the header line
	Lines  []*Line
}

type Line struct {
	Line int      // line number in the original input
	Kind LineKind // follows, goes to, fleet, land, scouts, status
	Text []byte   // text of the move
}

// warning - lines are sorted by kind first, then by line number,
// so keep that in mind when updating the LinkKind enum.

type LineKind int

const (
	Unknown LineKind = iota
	Turn
	UnitFollows
	UnitGoesTo
	UnitLandMove
	UnitFleetMove
	UnitScouts
	UnitStatus
)

// Split splits the input report into sections.
// Each section contains the header and move data for a single unit.
// All other lines are ignored.
//
// We assume the caller has not done any clean up on the input.
//
// Returns a list of sections sorted by clan id, then by unit id, then by line number.
// We would like to include turn number in the sort, but we can't trust the report to include the turn number.
//
// Warnings:
//   - All input is converted to lowercase to make comparisons easier in future stages.
//   - Report sections sometimes are missing the Status line. We can't depend on it to close out a section.
//   - Sections can contain multiple turn lines because of the missing Status line. When that happens,
//     we capture the additional turn lines and hope that someone eventually reports an error.
func Split(input [][]byte) (sections []*Section) {
	var section *Section
	var prevClanId int
	for no, line := range input {
		//log.Printf("section: %d: %q\n", no, line)
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// force to lowercase and compress whitespace
		line = norm.NormalizeSpaces(bytes.ToLower(line))
		//log.Printf("section: %d: %q\n", no, line)
		// only look for a unit header if we are not in a section
		if is.UnitHeader(line) {
			if section != nil {
				section.Sort()
			}
			section = &Section{
				Id:     len(sections) + 1,
				Line:   no + 1,
				Header: bdup(line),
			}
			sections = append(sections, section)
			// hack the line apart and extract the unit id.
			// we know that the line looks like `tribe tribe_id,` because it is a Unit Header.
			unitUnitId, _, _ := bytes.Cut(line, []byte{','})
			_, section.UnitId, _ = bytes.Cut(unitUnitId, []byte{' '})
			section.ClanId, _ = strconv.Atoi(string(section.UnitId[1:4]))
			if section.ClanId != prevClanId {
				//log.Printf("section: %6d: %04d: unit %q\n", no+1, section.ClanId, section.UnitId)
				prevClanId = section.ClanId
			}
		} else if section == nil {
			//log.Printf("section: %d: ignoring line %q\n", no, line)
			continue
		} else if is.FleetMovement(line) {
			section.Lines = append(section.Lines, &Line{Line: no + 1, Kind: UnitFleetMove, Text: line})
		} else if is.TribeFollows(line) {
			section.Lines = append(section.Lines, &Line{Line: no + 1, Kind: UnitFollows, Text: line})
		} else if is.TribeGoesTo(line) {
			section.Lines = append(section.Lines, &Line{Line: no + 1, Kind: UnitGoesTo, Text: line})
		} else if is.TribeMovement(line) {
			section.Lines = append(section.Lines, &Line{Line: no + 1, Kind: UnitLandMove, Text: line})
		} else if is.ScoutLine(line) {
			section.Lines = append(section.Lines, &Line{Line: no + 1, Kind: UnitScouts, Text: line})
		} else if is.TurnHeader(line) {
			section.Lines = append(section.Lines, &Line{Line: no + 1, Kind: Turn, Text: line})
		} else if is.UnitStatus(line) {
			section.Lines = append(section.Lines, &Line{Line: no + 1, Kind: UnitStatus, Text: line})
			// we would like to set `section` to nil to avoid capturing lines between sections,
			// but turn reports sometimes don't include the Status line. that means that we are
			// sometimes going to include movement lines that belong to a different unit.
			// bytes.HasPrefix(line, section.UnitId) would help us catch that condition.
		}
	}
	if section != nil {
		section.Sort()
	}

	// we would like to sort the sections by unit then turn, but we can't trust the turn data in any report.
	// so we just sort by clan, then unit, then line number.
	sort.Slice(sections, func(i, j int) bool {
		if sections[i].ClanId < sections[j].ClanId {
			return true
		} else if sections[i].ClanId == sections[j].ClanId {
			sop := bytes.Compare(sections[i].UnitId, sections[j].UnitId)
			if sop < 0 {
				return true
			} else if sop == 0 {
				return sections[i].Line < sections[j].Line
			}
		}
		return false
	})

	return sections
}

// bdup returns a copy of the slice.
func bdup(b []byte) []byte {
	b2 := make([]byte, len(b))
	copy(b2, b)
	return b2
}

func (s *Section) Sort() {
	// sort by kind, then by line number
	sort.Slice(s.Lines, func(i, j int) bool {
		if s.Lines[i].Kind != s.Lines[j].Kind {
			return s.Lines[i].Kind < s.Lines[j].Kind
		}
		return s.Lines[i].Line < s.Lines[j].Line
	})
}

func DumpSections(sections []*Section) []byte {
	b := &bytes.Buffer{}
	for _, s := range sections {
		b.Write(s.Header)
		b.WriteByte('\n')
		for _, line := range s.Lines {
			b.Write(line.Text)
			b.WriteByte('\n')
		}
	}
	return b.Bytes()
}
