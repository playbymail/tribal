// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

import (
	"unicode/utf8"
)

// RemoveBadUtf8 processes a byte slice and replaces any invalid UTF-8 sequences
// with the UTF-8 replacement character.
// Returns a new byte slice containing only the valid UTF-8 sequences.
func RemoveBadUtf8(input []byte) []byte {
	const runeErrorByte = "\uFFFD" // UTF-8 replacement character

	if input == nil {
		return []byte{}
	}
	output := make([]byte, 0, len(input))
	for i := 0; i < len(input); {
		// if the encoding is invalid, DecodeRune returns (RuneError, 1).
		// otherwise, it returns (r, w) where r is the rune and w is the width of the run, in bytes.
		r, w := utf8.DecodeRune(input)
		if r == utf8.RuneError {
			// invalid sequence; copy replacement character
			output = append(output, runeErrorByte...)
		} else {
			// valid rune; copy it
			output = append(output, input[i:i+w]...)
		}
		// advance input by the width of the rune we just processed
		i += w
	}

	return output
}
