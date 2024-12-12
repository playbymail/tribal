// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package stdlib

import "bytes"

// dedup removes duplicate elements from a slice.
// it assumes that the slice is sorted.
func dedup(input [][]byte) [][]byte {
	if len(input) == 0 {
		return input
	}
	var lines [][]byte
	var prevLine []byte
	// remove duplicate lines from the sorted input
	for _, line := range input {
		if bytes.Compare(line, prevLine) != 0 {
			lines = append(lines, line)
			prevLine = line
		}
	}
	return lines
}
