// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

// TrimBlankLines removes leading and trailing blank lines from the slice of byte slices.
// Returns the trimmed slice.
// If the input slice is empty or contains only blank lines, returns nil.
func TrimBlankLines(lines [][]byte) [][]byte {
	if lines == nil {
		return nil
	}
	start, end := 0, len(lines)
	// find the first non-blank line
	for start < end && len(lines[start]) == 0 {
		start++
	}
	// find the last non-blank line
	for start < end && len(lines[end-1]) == 0 {
		end--
	}
	if start == end {
		return nil
	}
	return lines[start:end]
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
