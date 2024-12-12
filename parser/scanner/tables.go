// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

var (
	delimiters [256]bool
	glyphs     [256]bool
)

func init() {
	// digits and letters
	for _, ch := range []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`) {
		glyphs[ch] = true
	}
	// punctuation
	for _, ch := range []byte(`-/$.()>#,:_`) {
		glyphs[ch], delimiters[ch] = true, true
	}
	glyphs['\\'], delimiters['\\'] = true, true
	// newline
	glyphs['\n'], delimiters['\n'] = true, true
	// quote marks
	glyphs['"'], delimiters['"'] = true, true
	glyphs['\''], delimiters['\''] = true, true
}
