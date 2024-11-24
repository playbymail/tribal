// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

import (
	"bytes"
	"github.com/playbymail/tribal/is"
)

// MappingLines returns a copy of the lines in the input that match the following:
// - Unit headers
// - Turn headers
// - Movement lines
// - Unit status lines
//
// Assumes that the input has already been cleaned up and converted to lower case.
func MappingLines(input []byte) [][]byte {
	lines := bytes.Split(input, []byte{'\n'})
	output := make([][]byte, 0, len(lines))
	for _, line := range lines {
		if is.UnitHeader(line) || is.TurnHeader(line) || is.MovementLine(line) || is.UnitStatus(line) {
			dup := make([]byte, len(line))
			copy(dup, line)
			output = append(output, dup)
		}
	}
	return output
}
