// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

const (
	// ASCII codes for carriage return and line feed.
	CR = '\r'
	LF = '\n'
)

// LineEndings returns a copy of the input with all line endings
// (Windows, MacOS, or Unix) converted to Unix EOL.
func LineEndings(data []byte) []byte {
	// buffer to store normalized output.
	normalized := make([]byte, 0, len(data))
	// process each byte in the input slice.
	for i := 0; i < len(data); i++ {
		if data[i] != CR {
			// copy non-CR characters as-is.
			normalized = append(normalized, data[i])
			continue
		}
		// handle Windows-style CR LF or standalone CR.
		if i+1 < len(data) && data[i+1] == LF {
			// skip the CR in CR LF.
			i++
		}
		// replace with a single LF.
		normalized = append(normalized, LF)
	}
	return normalized
}
