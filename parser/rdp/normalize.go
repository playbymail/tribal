// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package rdp

const (
	// CR and LF are ASCII codes for carriage return and line feed.
	CR = '\r'
	LF = '\n'
)

var (
	// pre-computed lookup table for delimiters
	isSpaceDelimiter [256]bool
)

func init() {
	// initialize the lookup table for delimiters
	for _, ch := range []byte{CR, LF, ',', '(', ')', '\\', ':'} {
		isSpaceDelimiter[ch] = true
	}
}

// normalize returns a copy of the input with:
//  1. all runs of spaces and tabs replaced by a single space
//  2. insignificant spaces (e.g., before and after delimiters) removed
//  3. upper-case ASCII converted to lower-case
//  4. Windows and MacOS line endings converted to Unix eol
func normalize(input []byte) []byte {
	if len(input) == 0 {
		return []byte{}
	}

	// buffer to store normalized output.
	normalized := make([]byte, 0, len(input))

	prevCharWasDelimiter := false

	for i := 0; i < len(input); i++ {
		if isSpaceOrTab(input[i]) { // skip runs of spaces and tabs
			for ; i < len(input) && isSpaceOrTab(input[i]); i++ {
				// skip all spaces/tabs
			}
			// check if the space is significant or not
			nextCharIsDelimiter := i >= len(input) || isSpaceDelimiter[input[i]]
			isSignificantSpace := !(prevCharWasDelimiter || nextCharIsDelimiter)
			if isSignificantSpace {
				// space is significant, so keep it
				normalized = append(normalized, ' ')
			}
			i-- // adjust for the outer loop increment
			continue
		} else if 'A' <= input[i] && input[i] <= 'Z' {
			// write the current character and update the delimiter state
			normalized = append(normalized, input[i]-'A'+'a')
			prevCharWasDelimiter = false
			continue
		} else if input[i] != CR {
			// write the current character and update the delimiter state
			normalized = append(normalized, input[i])
			prevCharWasDelimiter = isSpaceDelimiter[input[i]]
			continue
		}
		// handle Windows-style CR LF or standalone CR.
		if i+1 < len(input) && input[i+1] == LF {
			// skip the CR in CR LF.
			i++
		}
		// replace with a single LF and update the delimiter state
		normalized = append(normalized, LF)
		prevCharWasDelimiter = true
	}

	return normalized
}

// helper function to identify spaces and tabs
func isSpaceOrTab(b byte) bool {
	return b == ' ' || b == '\t'
}
