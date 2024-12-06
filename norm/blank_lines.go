// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

// RemoveEmptyLines removes all empty lines from the slice of byte slices.
// Returns a copy of the slice containing the lines that were not empty.
// If the input slice is empty or contains only empty lines, returns nil.
func RemoveEmptyLines(input [][]byte) [][]byte {
	if input == nil {
		return nil
	}
	var lines [][]byte
	for n, line := range input {
		if len(line) > 0 {
			if lines == nil {
				lines = make([][]byte, 0, len(input)-n)
			}
			lines = append(lines, line)
		}
	}
	return lines
}

// TrimLeadingBlankLines trims the leading blank lines from the slice of byte slices.
// Returns the trimmed slice. If the input slice is empty or contains only blank lines,
// returns an empty slice.
func TrimLeadingBlankLines(lines [][]byte) [][]byte {
	if lines == nil {
		return nil
	}
	start := 0
	for start < len(lines) && len(lines[start]) == 0 {
		start++
	}
	return lines[start:]
}

// TrimTrailingBlankLines trims the trailing blank lines from the slice of byte slices.
// Returns the trimmed slice. If the input contains only blank lines, returns an empty slice.
func TrimTrailingBlankLines(lines [][]byte) [][]byte {
	if lines == nil {
		return nil
	}
	end := len(lines)
	for end > 0 && len(lines[end-1]) == 0 {
		end--
	}
	return lines[:end]
}
