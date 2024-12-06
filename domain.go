// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package tribal

// this file defines the domain model for the application

// ClanId_t is the unique identifier for a clan.
// The clan id must be between 1 and 999.
// Clan's own reports and units. All queries should be filtered by clan id.
// As a special case, clan 188 is the GM's clan.
type ClanId_t int

// TurnId_t is the unique identifier for a turn.
// The range is 0 ... 9999 and starts at 0 for turn 899-12.
type TurnId_t int // turn number of the report

// UnitId_t is the unique identifier for a unit.
// It matches the pattern of type followed by an optional code and sequence number.
type UnitId_t string

// UnitName_t is the domain model for a unit name.
type UnitName_t string

// ReportFile_t is the domain model for a report file.
// All reports are owned by a clan and only that clan should have access to the report.
//
// The report file is the source of truth for the report, but we don't store the original report.
// Instead, we store the normalized contents of the report because it is easier to work with.
//
// There are some errors that will prevent a report from being imported into the database.
// All other errors are stored in the database with the report.
//
// Report turn, name, and hash are unique within a clan.
//
// It's probably rude of us, but bits of the report that have errors are not properly rendered.
// For example, if there's an error with the turn number for a unit, the render will skip that unit.
type ReportFile_t struct {
	Owner ClanId_t  // clan that owns the report
	Name  string    // name of the report file (includes the path and is unique within a clan)
	Turn  TurnId_t  // turn number of the report (unique within a clan)
	Units []*Unit_t // all units in the report, in the order they appear in the report
	Error error     // highest level error encountered while parsing the report
	Hash  string    // SHA1 hash of the report's original contents (unique within a clan)
	Lines string    // report lines after normalization
}

// Turn_t is the domain model for a turn.
// Turn data is shared by all clans.
type Turn_t struct {
	Id    TurnId_t // unique identifier for the turn
	Year  int      // year of the turn, range is 899 ... 9999
	Month int      // month of the turn, range is 1 ... 12 (but year 899 only has month 12)
	Error error    // highest level error encountered while parsing the turn
}

// Unit_t is the domain model for a unit.
// Unit data is owned by a clan and never shared.
type Unit_t struct {
	Id    UnitId_t   // unique identifier for the unit
	Name  UnitName_t // optional name of the unit. mostly ignored.
	Error error      // highest level error encountered while parsing the unit
}

func (t TurnId_t) YearMonth() (int, int) {
	offset := int(t) + 12 // push 899-12 to end of year
	years := offset / 12
	months := offset - years*12
	if months == 0 { // trust me, this is correct
		years, months = years-1, 12
	}
	return years + 899, months
}
