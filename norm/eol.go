// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

import "bytes"

const (
	// CR and LF are the ASCII codes for carriage return and line feed, respectively.
	// They are used to represent the end of a line in text files and are needed for
	// cleaning up the text from Windows and MacOS line endings.
	CR = '\r'
	LF = '\n'
)

// NormalizeEOL returns a copy of the input with different types of EOL converted to Unix EOL.
// Converts Windows EOL (CR+LF) to Unix EOL (LF).
// Converts Classic Mac EOL (CR) to Unix EOL (LF).
// Unix EOL (LF) passes through unchanged.
//
// BUG: if the input is empty, it returns the input, not a copy of it.
func NormalizeEOL(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	output := bytes.NewBuffer(make([]byte, 0, len(input)))
	for len(input) != 0 {
		if input[0] == CR { // window or maybe classic mac
			input = input[1:]
			// found CR, check for CR LF
			if len(input) != 0 && input[0] == LF {
				input = input[1:]
			}
			output.WriteByte(LF)
			continue
		}
		output.WriteByte(input[0])
		input = input[1:]
	}
	return output.Bytes()
}
