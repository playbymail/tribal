// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package parser

//go:generate pigeon -o grammar.go grammar.peg

import (
	"crypto/sha1"
	"fmt"
	"github.com/playbymail/tribal/docx"
	"github.com/playbymail/tribal/norm"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/parser/turns"
	"github.com/playbymail/tribal/parser/units"
	"github.com/playbymail/tribal/text"
	"log"
)

// Report parses the report and returns the report.
// Tries to parse as much as possible before giving up and returning an error.
// Path is the path to the report. It's only used for logging.
// The parser is responsible for normalizing the input (line endings, case, etc.)
func Report(path string, options ...Option_t) (*Report_t, error) {
	rpt := &Report_t{}

	for _, opt := range options {
		if err := opt(&rpt.options); err != nil {
			rpt.Error = err
			return rpt, err
		}
	}
	if rpt.options.Data == nil {
		rpt.Error = ErrNoData
		return rpt, rpt.Error
	}
	rpt.Hash = rpt.options.Hash
	if rpt.Hash == "" {
		rpt.Hash = fmt.Sprintf("%x", sha1.Sum(rpt.options.Data))
	}

	// guess the file type and then extract the report into a slice of lines
	var lines [][]byte
	switch docx.DetectWordDocType(rpt.options.Data) {
	case docx.Doc:
		rpt.Error = ErrInvalidFileType
		return rpt, rpt.Error
	case docx.Docx:
		if wordLines, err := docx.Read(rpt.options.Data); err != nil {
			rpt.Error = err
			return rpt, err
		} else {
			lines = wordLines
		}
	default:
		if textLines, err := text.Read(rpt.options.Data); err != nil {
			rpt.Error = err
			return rpt, err
		} else {
			lines = textLines
		}
	}

	// normalize the input. this is what we will save to the database.
	// there are some users that can't edit Word documents.
	// maybe we'll allow them to use this file.
	// skips any empty lines.
	for n, line := range lines {
		// normalize the input
		line = norm.NormalizeSpaces(norm.NormalizeCase(norm.RemoveBadUtf8(line)))
		if len(line) == 0 {
			continue
		}
		if rpt.Lines == nil {
			// cheap attempt to optimize memory allocation
			rpt.Lines = make([][]byte, 0, len(lines)-n)
		}
		rpt.Lines = append(rpt.Lines, line)
	}

	// highest level loop in the parser looks for unit headings.
	// every time we find one, we pass control to a lower level of the parser.
	// that level returns the number of lines parsed and the highest level error it encountered.
	var currentUnit *ast.UnitHeading_t
	for _, line := range rpt.Lines {
		// is this line a unit heading?
		if uht := units.ParseUnitHeading(path, line); uht != nil {
			currentUnit = uht
			unit := &Unit_t{Id: UnitId_t(uht.Id)}
			rpt.Units = append(rpt.Units, unit)
			// log.Printf("import: report: unit header: %+v\n", *uht)
			continue
		}
		// is this line a turn number?
		if turnLine := turns.ParseTurnLine(path, line); turnLine != nil {
			rpt.Turn = &Turn_t{
				No:    turnLine.No,
				Year:  turnLine.Year,
				Month: turnLine.Month,
				Error: turnLine.Error,
			}
			// log.Printf("import: report: turn number: %+v\n", *turnLine)
			continue
		}
		// is this line a tribe movement?
		// is this line a tribe follows?
		// is this line a tribe goes to?
		// is this line a fleet movement?
		// is this line a scouting report?
		// is this line a tribe status?
	}
	if rpt.Turn == nil {
		// reject the entire report if there is no turn number found in the report.
		return nil, ErrNoTurnNumber
	} else if currentUnit == nil {
		// reject the entire report if there were not units found in the report.
		return nil, ErrNoUnits
	}

	// all units in the report must have the same turn number.
	foundUnexpectedTurnNumber := false
	for _, unit := range rpt.Units {
		if unit.Turn == nil {
			// patch the turn number for any unit that does not have a turn number.
			unit.Turn = rpt.Turn
			continue
		} else if unit.Turn.No == rpt.Turn.No {
			// this is good
			continue
		}
		log.Printf("import: report: unit %s: turn number: want %d, got %d\n", unit.Id, rpt.Turn.No, unit.Turn.No)
		unit.Error = ErrUnexpectedTurnNumber
		foundUnexpectedTurnNumber = true
	}
	if foundUnexpectedTurnNumber {
		rpt.Error = ErrUnexpectedTurnNumber
		return rpt, ErrUnexpectedTurnNumber
	}

	//// split the report into sections
	//ss := section.Split(lines)
	//log.Printf("import: report: %s: %d sections\n", path, len(ss))
	//
	//// reject if there are no sections
	//if len(ss) == 0 {
	//	rpt.Error = ErrNoSections
	//	return rpt, rpt.Error
	//}
	//// reject if the first section does not include the turn number
	//if ss[0].Turn == nil {
	//	rpt.Error = ErrMissingTurnNumberInFirstSection
	//	return rpt, rpt.Error
	//}
	//
	//// now parse each section
	//// create a variable to hold the parsed section (ast, cst?)
	//for _, us := range ss {
	//	sct := &Section_t{}
	//	rpt.Sections = append(rpt.Sections, sct)
	//
	//	log.Printf("import: report: %s: %3d: %5d: %s\n", path, us.Id, us.Line, us.Header)
	//	//// unit header compares only the first field (the unit type and id).
	//	//// we have to split the line to get the current hex.
	//	//for i, f := range bytes.Split(line, []byte{','}) {
	//	//	f = bytes.TrimSpace(f)
	//	//	if i == 0 {
	//	//		// we know that the first field is the unit type and id.
	//	//		section.Header.Unit = bdup(bytes.Fields(f)[1])
	//	//	} else if bytes.HasPrefix(f, []byte("current hex = ")) {
	//	//		section.Header.CurrentHex = bdup(f[len("current hex = "):])
	//	//	}
	//	//}
	//	//if section.Header.Unit == nil {
	//	//	// absolutely must have a unit id to be a unit header.
	//	//	continue
	//	//}
	//	//if section.Header.CurrentHex == nil {
	//	//	// badly broken report, so use N/A for the current hex.
	//	//	// this should cause the unit to be placed on a random hex
	//	//	// in the map with an appropriate error note.
	//	//	section.Header.CurrentHex = bdup([]byte{'n', '/', 'a'})
	//	//}
	//	hdr := Section_t{
	//		Unit:  UnitId_t(us.Header.Unit),
	//		Error: nil,
	//	}
	//	log.Printf("import: report: %s: %3d: %5d: %+v\n", path, us.Id, us.Line, hdr)
	//	if us.Id == 1 {
	//		log.Printf("import: report: %s: %3d: %5d: %s\n", path, us.Id, us.Line, us.Turn)
	//	}
	//	//if len(us.Moves.Movement) > 0 {
	//	//	us.Moves.Movement = preProcessTribeMovement(us.Moves.Movement)
	//	//	log.Printf("import: report: %s: %3d: %5d: %s\n", path, us.Id, us.Line, us.Moves.Movement)
	//	//}
	//	//if len(us.Moves.Fleet) > 0 {
	//	//	us.Moves.Fleet = preProcessFleetMovement(us.Moves.Fleet)
	//	//	log.Printf("import: report: %s: %3d: %5d: %q\n", path, us.Id, us.Line, us.Moves.Fleet)
	//	//}
	//	//for i := range us.Moves.Scouts {
	//	//	if len(us.Moves.Scouts[i]) == 0 {
	//	//		continue
	//	//	}
	//	//	us.Moves.Scouts[i] = preProcessScoutMovement(us.Moves.Scouts[i])
	//	//	log.Printf("import: report: %s: %3d: %5d: scout %s\n", path, us.Id, us.Line, us.Moves.Scouts[i])
	//	//}
	//	//us.Status = preProcessUnitStatus(us.Status)
	//	//log.Printf("import: report: %s: %3d: %5d: %s\n", path, us.Id, us.Line, us.Status)
	//}

	return rpt, nil
}

func parseSectionHeader(path string, line []byte, debug bool) (SectionHeader_t, error) {
	if va, err := Parse(path, line, Entrypoint("SectionHeader")); err != nil {
		log.Printf("%s: section header: %s: %v\n", path, slug(line, 14), err)
		return SectionHeader_t{}, err
	} else if sectionHeader, ok := va.(SectionHeader_t); !ok {
		log.Printf("%s: section header: %q\n", path, slug(line, 15))
		log.Printf("error: invalid type\n")
		log.Printf("please report this error")
		panic(fmt.Errorf("want SectionHeader_t, got %T", va))
	} else {
		return sectionHeader, nil
	}
}

func slug(b []byte, n int) string {
	if len(b) < n {
		return string(b)
	}
	return string(b[:n])
}
