// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package section

import (
	"bytes"
	"fmt"
	"github.com/playbymail/tribal/domains"
	"github.com/playbymail/tribal/section/status"
	"github.com/playbymail/tribal/section/turns"
	"github.com/playbymail/tribal/section/units"
	"log"
)

//go:generate pigeon -o turns/grammar.go turns/grammar.peg
//go:generate pigeon -o units/grammar.go units/grammar.peg

type Section struct {
	Line   int    // line number in the original input
	Id     int    // section number, starting at 1
	ClanId int    // derived from the header
	UnitId []byte // taken from the header
	Lines  struct {
		FleetMoves  []byte
		Turn        []byte
		ScoutLines  [][]byte
		Status      []byte
		Unit        []byte // text of the header line
		UnitFollows []byte
		UnitGoesTo  []byte
		UnitMoves   []byte
	}
	Unit   *domains.Unit_t
	Errors []error // error from parsing the unit header
}

// Less returns true if section should be sorted before another section.
// We sort by clan, then unit, then line number.
func (s *Section) Less(s2 *Section) bool {
	if s.ClanId < s.ClanId {
		return true
	} else if s.ClanId == s.ClanId {
		sop := bytes.Compare(s.UnitId, s.UnitId)
		if sop < 0 {
			return true
		} else if sop == 0 {
			return s.Line < s.Line
		}
	}
	return false
}

func (s *Section) Parse(path string) error {
	// sort the lines before parsing.
	s.Sort()

	var ok bool

	// the header is the only mandatory line
	if s.Lines.Unit == nil {
		return fmt.Errorf("section %d: missing unit header", s.Id)
	}
	if v, err := units.Parse(path, s.Lines.Unit); err != nil {
		s.Errors = append(s.Errors, err)
		log.Printf("section: header %q: parse error %v\n", s.Lines.Unit, err)
		return err
	} else if s.Unit, ok = v.(*domains.Unit_t); !ok {
		panic(fmt.Sprintf("assert(%T == *UnitHeading_t)", v))
	} else {
		//unitHeading.Name = uh.Name
		//unitHeading.CurrentHex = uh.CurrentHex
		//unitHeading.PreviousHex = uh.PreviousHex
		//unitHeading.Error = uh.Error
		log.Printf("section: header %q: unit heading: %+v", s.Lines.Unit, *s.Unit)
	}

	// turn line is optional, even though it's required in the spec
	if s.Lines.Turn != nil {
		if v, err := turns.Parse(path, s.Lines.Turn); err != nil {
			s.Errors = append(s.Errors, err)
			log.Printf("section: turn %q: parse error %v\n", s.Lines.Turn, err)
		} else if s.Unit.Turn, ok = v.(domains.TurnId_t); !ok {
			panic(fmt.Sprintf("assert(%T == TurnId_t)", v))
		} else {
			log.Printf("section: turn %q: %+v", s.Lines.Turn, s.Unit.Turn)
		}
	}

	// parse movement line. note that we will never parse more than one movement line.
	// we arbitrarily choose to test in order of follows, goes to, unit and then fleet moves.

	// next lines should be the scouting lines

	// status line is optional, even though it's required in the spec.
	// if present, it must start with the unit id.
	if s.Lines.Status == nil {
		// should be an error but the setup reports often don't include it.
	} else if !bytes.HasPrefix(s.Lines.Status, []byte(fmt.Sprintf("%s status:", s.Unit.Id))) {
		s.Errors = append(s.Errors, domains.ErrInvalidStatusPrefix)
		log.Printf("section: status %q: invalid prefix", s.Lines.Status)
	} else {
		var err error
		if s.Unit.Status, err = status.Parse(path, s.Lines.Status); err != nil {
			s.Errors = append(s.Errors, err)
			log.Printf("section: status %q: parse error %v\n", s.Lines.Status, err)
		} else {
			log.Printf("section: status %q: %+v", s.Lines.Status, s.Unit.Status)
		}
	}

	log.Printf("section: %q\n\t%+v", s.Lines.Unit, *s.Unit)

	return nil
}

func (s *Section) Sort() {
	//// sort by kind, then by line number
	//sort.Slice(s.Lines, func(i, j int) bool {
	//	return s.Lines[i].Less(s.Lines[j])
	//})
}

func DumpSections(sections []*Section, separateUnits bool) []byte {
	b := &bytes.Buffer{}
	for n, s := range sections {
		if n > 0 && separateUnits {
			b.WriteByte('\n')
		}
		b.Write(s.Lines.Unit)
		b.WriteByte('\n')
		if len(s.Lines.Turn) != 0 {
			b.Write(s.Lines.Turn)
			b.WriteByte('\n')
		}
		for _, line := range [][]byte{s.Lines.UnitFollows, s.Lines.UnitGoesTo, s.Lines.UnitMoves, s.Lines.FleetMoves} {
			if len(line) != 0 {
				b.Write(line)
				b.WriteByte('\n')
			}
		}
		for _, line := range s.Lines.ScoutLines {
			if len(line) != 0 {
				b.Write(line)
				b.WriteByte('\n')
			}
		}
		if len(s.Lines.Status) != 0 {
			b.Write(s.Lines.Status)
			b.WriteByte('\n')
		}
	}
	return b.Bytes()
}
