// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package is

import (
	"bytes"
	"regexp"
)

var (
	rxCourierHeader  = regexp.MustCompile(`^courier \d{4}c\d,`)
	rxElementHeader  = regexp.MustCompile(`^element \d{4}e\d,`)
	rxFleetHeader    = regexp.MustCompile(`^fleet \d{4}f\d,`)
	rxGarrisonHeader = regexp.MustCompile(`^garrison \d{4}g\d,`)
	rxTribeHeader    = regexp.MustCompile(`^tribe \d{4},`)

	rxTurnHeader = regexp.MustCompile(`^current turn \d{3,4}-\d{1,2}\(#\d+\),`)

	rxFleetMovement = regexp.MustCompile(`^(calm|mild|strong|gale) (ne|se|sw|nw|n|s) fleet movement:`)
	rxScoutLine     = regexp.MustCompile(`^scout [1-8]:`)

	rxCourierStatus  = regexp.MustCompile(`^\d{4}c\d status:`)
	rxElementStatus  = regexp.MustCompile(`^\d{4}e\d status:`)
	rxFleetStatus    = regexp.MustCompile(`^\d{4}f\d status:`)
	rxGarrisonStatus = regexp.MustCompile(`^\d{4}g\d status:`)
	rxTribeStatus    = regexp.MustCompile(`^\d{4} status:`)
)

// FleetMovement returns true if the line represents a fleet movement.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func FleetMovement(line []byte) bool {
	return rxFleetMovement.Match(line)
}

// MovementLine returns true if the line represents a unit movement line.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func MovementLine(line []byte) bool {
	return TribeMovement(line) || TribeFollows(line) || TribeGoesTo(line) || ScoutLine(line) || FleetMovement(line)
}

// ScoutLine determines if a line represents a TribeNet scout command.
// Example: "scout 1: scout s-pr"
//
// Assumes that the line has already been cleaned up and converted to lower case.
func ScoutLine(line []byte) bool {
	return rxScoutLine.Match(line)
}

// TribeFollows returns true if the line represents a tribe follows command.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func TribeFollows(line []byte) bool {
	return bytes.HasPrefix(line, []byte("tribe follows "))
}

// TribeGoesTo returns true if the line represents a tribe goes to command.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func TribeGoesTo(line []byte) bool {
	return bytes.HasPrefix(line, []byte("tribe goes to "))
}

// TribeMovement returns true if the line represents a tribe movement line.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func TribeMovement(line []byte) bool {
	return bytes.HasPrefix(line, []byte("tribe movement:"))
}

// TurnHeader returns true if a line represents a TribeNet turn header.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func TurnHeader(line []byte) bool {
	return rxTurnHeader.Match(line)
}

// UnitHeader returns true if a line represents a TribeNet unit header.
// It checks for five different types of unit headers:
//   - Tribe headers
//   - Courier headers
//   - Element headers
//   - Fleet headers
//   - Garrison headers
//
// Returns true if the line matches any of these header patterns.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func UnitHeader(line []byte) bool {
	return rxTribeHeader.Match(line) || rxCourierHeader.Match(line) || rxElementHeader.Match(line) || rxFleetHeader.Match(line) || rxGarrisonHeader.Match(line)
}

// UnitStatus returns true if a line represents a TribeNet unit status line.
// It checks for five different types of unit status lines:
//   - Tribe status
//   - Courier status
//   - Element status
//   - Fleet status
//   - Garrison status
//
// Returns true if the line matches any of these status line patterns.
//
// Assumes that the line has already been cleaned up and converted to lower case.
func UnitStatus(line []byte) bool {
	return rxTribeStatus.Match(line) || rxCourierStatus.Match(line) || rxElementStatus.Match(line) || rxFleetStatus.Match(line) || rxGarrisonStatus.Match(line)
}
