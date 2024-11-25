// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

var (
	// pre-computed lookup table for acceptable printing characters
	isPrintingGlyph [256]bool
)

func init() {
	// initialize the lookup table for acceptable printing characters
	for ch := '!'; ch < '~'; ch++ {
		isPrintingGlyph[ch] = true
	}
	isPrintingGlyph['\n'] = true
}

// PrintingGlyphs returns the slice with all non-printing characters
// replace with spaces. Updates the input slice in-place.
func PrintingGlyphs(input []byte) []byte {
	for i := 0; i < len(input); i++ {
		if !isPrintingGlyph[input[i]] {
			input[i] = ' '
		}
	}
	return input
}
