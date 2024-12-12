// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

import (
	"bytes"
	"regexp"
)

var (
	reBackslashDash = regexp.MustCompile(`\\+-+ *`)

	reBackslashComma = regexp.MustCompile(`\\+,+`)
	reBackslashUnit  = regexp.MustCompile(`\\+(\d{4}(?:[cefg]\d)?)`)

	reCommaBackslash         = regexp.MustCompile(`,+\\`)
	reCommaCanal             = regexp.MustCompile(`,canal `)
	reCommaCantMove          = regexp.MustCompile(`,can't move `)
	reCommaFind              = regexp.MustCompile(`,find `)
	reCommaFord              = regexp.MustCompile(`,ford `)
	reCommaHsm               = regexp.MustCompile(`,hsm `)
	reCommaL                 = regexp.MustCompile(`,l `)
	reCommaLcm               = regexp.MustCompile(`,lcm `)
	reCommaLjm               = regexp.MustCompile(`,ljm `)
	reCommaLsm               = regexp.MustCompile(`,lsm `)
	reCommaNoFord            = regexp.MustCompile(`,no ford `)
	reCommaNotEnoughMPs      = regexp.MustCompile(`,not enough m.p's `)
	reCommaNothingOfInterest = regexp.MustCompile(`,nothing of interest `)
	reCommaO                 = regexp.MustCompile(`,o `)
	reCommaPass              = regexp.MustCompile(`,pass `)
	reCommaPatrolledAndFound = regexp.MustCompile(`,patrolled and found `)
	reCommaRiver             = regexp.MustCompile(`,river `)
	reCommaStoneRoad         = regexp.MustCompile(`,stone road `)

	reDirectionCommaUnit = regexp.MustCompile(`\b(ne|se|sw|nw|n|s)[ ,](\d{4}(?:[cefg]\d)?)`)

	reDirectionUnit = regexp.MustCompile(`\b(ne|se|sw|nw|n|s) (\d{4}(?:[cefg]\d)?)`)

	reRunOfBackslashes = regexp.MustCompile(`\\\\+`)
	reRunOfComma       = regexp.MustCompile(`,,+`)
)

// preProcessCurrentTurn removes the optional turn data from the line.
func preProcessCurrentTurn(line []byte) []byte {
	turnInfo, _, ok := bytes.Cut(line, []byte{','})
	if !ok {
		return line
	}
	return turnInfo
}

// preProcessFleetMovement processes a fleet movement line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func preProcessFleetMovement(line []byte) []byte {
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

// preProcessScoutMovement processes a scout movement line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func preProcessScoutMovement(line []byte) []byte {
	// replace some common comma patterns with semicolons
	line = reCommaCanal.ReplaceAll(line, []byte(";canal "))
	line = reCommaCantMove.ReplaceAll(line, []byte(";can't move "))
	line = reCommaFind.ReplaceAll(line, []byte(";find "))
	line = reCommaFord.ReplaceAll(line, []byte(";ford "))
	line = reCommaHsm.ReplaceAll(line, []byte(";hsm "))
	line = reCommaL.ReplaceAll(line, []byte(";l "))
	line = reCommaLcm.ReplaceAll(line, []byte(";lcm "))
	line = reCommaLjm.ReplaceAll(line, []byte(";ljm "))
	line = reCommaLsm.ReplaceAll(line, []byte(";lsm "))
	line = reCommaNoFord.ReplaceAll(line, []byte(";no ford "))
	line = reCommaNotEnoughMPs.ReplaceAll(line, []byte(";not enough m.p's "))
	line = reCommaNothingOfInterest.ReplaceAll(line, []byte(";nothing of interest "))
	line = reCommaO.ReplaceAll(line, []byte(";o "))
	line = reCommaPass.ReplaceAll(line, []byte(";pass "))
	line = reCommaPatrolledAndFound.ReplaceAll(line, []byte(";patrolled and found "))
	line = reCommaRiver.ReplaceAll(line, []byte(";river "))
	line = reCommaStoneRoad.ReplaceAll(line, []byte(";stone road "))

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

	return line
}

// preProcessSpecialHex is a no-op.
func preProcessSpecialHex(line []byte) []byte {
	return line
}

// preProcessTribeFollows is a no-op.
func preProcessTribeFollows(line []byte) []byte {
	return line
}

// preProcessTribeGoes is a no-op.
func preProcessTribeGoes(line []byte) []byte {
	return line
}

// preProcessTribeMovement processes a tribe movement line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func preProcessTribeMovement(line []byte) []byte {
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

	return line
}

// preProcessTurnLine is a no-op.
func preProcessTurnLine(line []byte) []byte {
	return line
}

// preProcessUnitStatus processes a unit status line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func preProcessUnitStatus(line []byte) []byte {
	// replace some common comma patterns with semicolons
	line = reCommaCanal.ReplaceAll(line, []byte(";canal "))
	line = reCommaFord.ReplaceAll(line, []byte(";ford "))
	line = reCommaHsm.ReplaceAll(line, []byte(";hsm "))
	line = reCommaL.ReplaceAll(line, []byte(";l "))
	line = reCommaLcm.ReplaceAll(line, []byte(";lcm "))
	line = reCommaLjm.ReplaceAll(line, []byte(";ljm "))
	line = reCommaLsm.ReplaceAll(line, []byte(";lsm "))
	line = reCommaO.ReplaceAll(line, []byte(";o "))
	line = reCommaPass.ReplaceAll(line, []byte(";pass "))
	line = reCommaRiver.ReplaceAll(line, []byte(";river "))
	line = reCommaStoneRoad.ReplaceAll(line, []byte(";stone road "))
	line = reDirectionCommaUnit.ReplaceAll(line, []byte{'$', '1', ';', '$', '2'})

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

	return line
}
