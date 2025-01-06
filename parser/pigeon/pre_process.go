// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package pigeon

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

	reRunOfBackslashes = regexp.MustCompile(`\\\\+`)
	reRunOfComma       = regexp.MustCompile(`,,+`)
)

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

// preProcessUnitStatus processes a unit status line to fix issues with backslash
// or direction followed by a unit ID. Caller must have already compressed spaces
// on the input line and forced to lowercase.
func preProcessUnitStatus(line []byte) []byte {
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
