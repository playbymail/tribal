// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

import "bytes"

// tokenizer is a simple tokenizer for TribeNet reports.
type tokenizer struct {
	input  []byte
	length int
	pos    int
}

// getch returns the next byte in the input and advances the position.
// Returns 0 if at end of input.
func (tk *tokenizer) getch() byte {
	if tk.pos >= tk.length {
		return 0
	}
	ch := tk.input[tk.pos]
	tk.pos++
	return ch
}

// match checks if the next byte matches any in the valid set.
// If there's a match, advances position and returns the byte with true.
// If no match or at end of input, returns 0 with false.
func (tk *tokenizer) match(valid ...byte) (byte, bool) {
	if tk.pos >= tk.length {
		return 0, false
	}
	ch := tk.input[tk.pos]
	if bytes.IndexByte(valid, ch) == -1 {
		return 0, false
	}
	tk.pos++
	return ch, true
}

// next returns the next token in the input.
// Handles newlines, whitespace, and text tokens distinctly.
// For text tokens, converts value to lowercase.
// Returns EOF token if at end of input.
func (tk *tokenizer) next() Token {
	if tk.pos >= tk.length {
		return Token{Type: EOF}
	}
	// anchor the token
	offset := tk.pos

	// get the first character of the token
	ch := tk.getch()
	if ch == '\n' {
		return Token{Type: Newline}
	}
	// lump whitespace and invalid characters together
	if !glyphs[ch] {
		for tk.pos < tk.length && !glyphs[tk.peek()] {
			_ = tk.getch()
		}
		return Token{Type: Whitespace}
	}
	// everything else is a text token
	if !delimiters[ch] {
		// if it's not a delimiter, keep reading until we hit one or whitespace
		for tk.pos < tk.length && glyphs[tk.peek()] && !delimiters[tk.peek()] {
			_ = tk.getch()
		}
	}
	// convert to lowercase and save the value as a string
	length := tk.pos - offset
	value := string(bytes.ToLower(tk.input[offset : offset+length]))
	// return it as a text token
	return Token{Type: Text, Value: value}
}

// peek returns the next byte in the input without advancing the position.
// Returns 0 if at end of input.
func (tk *tokenizer) peek() byte {
	if tk.pos >= tk.length {
		return 0
	}
	return tk.input[tk.pos]
}
