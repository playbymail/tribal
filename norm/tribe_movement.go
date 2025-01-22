// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package norm

import (
	"bytes"
	"regexp"
)

var (
	reBackslashDash = regexp.MustCompile(`\\+-+ *`)

	reBackslashComma = regexp.MustCompile(`\\+,+`)
	reBackslashUnit  = regexp.MustCompile(`\\+(\d{4}(?:[cefg]\d)?)`)
	reCommaBackslash = regexp.MustCompile(`,+\\`)
	reDirectionUnit  = regexp.MustCompile(`\b(ne|se|sw|nw|n|s) (\d{4}(?:[cefg]\d)?)`)

	// matches space direction comma
	reSpaceDirectionCommaDirection = regexp.MustCompile(` (nw|ne|n|sw|se|s),(?:nw|ne|n|sw|se|s)([,\\]|$)`)

	// matches a unit ID followed by comma followed by another unit ID
	reUnitCommaUnit = regexp.MustCompile(`([0-9]{4}(?:[cefg][1-9])?),([0-9]{4}(?:[cefg][1-9])?)`)

	reRunOfBackslashes = regexp.MustCompile(`\\\\+`)
	reRunOfComma       = regexp.MustCompile(`,,+`)
)

// FleetMovement processes a fleet movement line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func FleetMovement(line []byte) []byte {
	// replace backslash+dash with backslash
	line = reBackslashDash.ReplaceAll(line, []byte{'\\'})

	// replace backslash+comma and comma+backslash with backslash
	line = reBackslashComma.ReplaceAll(line, []byte{'\\'})
	line = reCommaBackslash.ReplaceAll(line, []byte{'\\'})

	// fix issues with backslash or direction followed by a unit ID
	line = reBackslashUnit.ReplaceAll(line, []byte{',', '$', '1'})
	line = reDirectionUnit.ReplaceAll(line, []byte{'$', '1', ',', '$', '2'})

	// reduce runs of certain punctuation to a single punctuation character
	line = reRunOfBackslashes.ReplaceAll(line, []byte{'\\'})
	line = reRunOfComma.ReplaceAll(line, []byte{','})

	// tweak the fleet movement to remove the trailing comma from the observations
	line = bytes.ReplaceAll(line, []byte{',', ')'}, []byte{')'})

	// remove all trailing backslashes from the line
	line = bytes.TrimRight(line, "\\")

	return line
}

// ScoutMovement processes a scout movement line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func ScoutMovement(line []byte) []byte {
	// replace backslash+dash with backslash
	line = reBackslashDash.ReplaceAll(line, []byte{'\\'})

	// replace backslash+comma and comma+backslash with backslash
	line = reBackslashComma.ReplaceAll(line, []byte{'\\'})
	line = reCommaBackslash.ReplaceAll(line, []byte{'\\'})

	// fix issues with backslash or direction followed by a unit ID
	line = reBackslashUnit.ReplaceAll(line, []byte{',', '$', '1'})
	line = reDirectionUnit.ReplaceAll(line, []byte{'$', '1', ',', '$', '2'})

	// reduce runs of certain punctuation to a single punctuation character
	line = reRunOfBackslashes.ReplaceAll(line, []byte{'\\'})
	line = reRunOfComma.ReplaceAll(line, []byte{','})

	// change the separator for nothing of interest
	if idx := bytes.Index(line, []byte(`,nothing of interest found`)); idx != -1 {
		line[idx] = '\\'
	}
	if idx := bytes.Index(line, []byte(`,patrolled and found `)); idx != -1 {
		line[idx] = '\\'
	}

	// remove all trailing backslashes from the line
	line = bytes.TrimRight(line, "\\")

	// cleanup lists of directions and units
	line = ListOfDirections(line)
	line = ListOfUnitIDs(line)

	return line
}

// TribeMovement processes a tribe movement line to fix issues with backslash or direction followed by a unit ID.
// Caller must have already compressed spaces on the input line and forced to lowercase.
func TribeMovement(line []byte) []byte {
	// force a backslash after the word moves
	if len(line) > 19 && line[19] == ' ' { // 	"tribe movement:move"
		line[19] = '\\'
	}

	// replace backslash+dash with backslash
	line = reBackslashDash.ReplaceAll(line, []byte{'\\'})

	// replace backslash+comma and comma+backslash with backslash
	line = reBackslashComma.ReplaceAll(line, []byte{'\\'})
	line = reCommaBackslash.ReplaceAll(line, []byte{'\\'})

	// fix issues with backslash or direction followed by a unit ID
	line = reBackslashUnit.ReplaceAll(line, []byte{',', '$', '1'})
	line = reDirectionUnit.ReplaceAll(line, []byte{'$', '1', ',', '$', '2'})

	// reduce runs of certain punctuation to a single punctuation character
	line = reRunOfBackslashes.ReplaceAll(line, []byte{'\\'})
	line = reRunOfComma.ReplaceAll(line, []byte{','})

	// remove all trailing backslashes from the line
	line = bytes.TrimRight(line, "\\")

	// cleanup lists of directions and units
	line = ListOfDirections(line)
	line = ListOfUnitIDs(line)

	return line
}

// UnitStatus processes a unit status line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func UnitStatus(line []byte) []byte {
	// replace backslash+dash with backslash
	line = reBackslashDash.ReplaceAll(line, []byte{'\\'})

	// replace backslash+comma and comma+backslash with backslash
	line = reBackslashComma.ReplaceAll(line, []byte{'\\'})
	line = reCommaBackslash.ReplaceAll(line, []byte{'\\'})

	// fix issues with backslash or direction followed by a unit ID
	line = reBackslashUnit.ReplaceAll(line, []byte{',', '$', '1'})
	line = reDirectionUnit.ReplaceAll(line, []byte{'$', '1', ',', '$', '2'})

	// reduce runs of certain punctuation to a single punctuation character
	line = reRunOfBackslashes.ReplaceAll(line, []byte{'\\'})
	line = reRunOfComma.ReplaceAll(line, []byte{','})

	// remove all trailing backslashes from the line
	line = bytes.TrimRight(line, "\\")

	// cleanup lists of directions and units
	line = ListOfDirections(line)
	line = ListOfUnitIDs(line)

	return line
}

// ListOfDirections replaces comma separated directions with space separated directions.
func ListOfDirections(line []byte) []byte {
	for {
		// FindIndex returns [start, end] of the first match
		loc := reSpaceDirectionCommaDirection.FindIndex(line)
		if loc == nil {
			break
		}
		// start at location and find the comma
		for i := loc[0]; i < loc[1]; i++ {
			if line[i] == ',' {
				// replace the comma with a space
				line[i] = ' '
				break
			}
		}
	}

	return line
}

// ListOfUnitIDs processes a line and replaces comma separated unit IDs with
// space separated unit IDs.
func ListOfUnitIDs(line []byte) []byte {
	for {
		// FindIndex returns [start, end] of the first match
		loc := reUnitCommaUnit.FindIndex(line)
		if loc == nil {
			break
		}
		// start at location and find the comma
		for i := loc[0]; i < loc[1]; i++ {
			if line[i] == ',' {
				// replace the comma with a space
				line[i] = ' '
				break
			}
		}
	}

	return line
}
