// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

import (
	"bytes"
)

// Tokenizer is the interface for tokenizing input.
// It should be implemented by any tokenizer and specialized to the version of the TribeNet report.
type Tokenizer interface {
	getch() byte
	match(...byte) (byte, bool)
	next() Token
	peek() byte
	peekNext() byte
}

// tokenizer_899_12 implements the Tokenizer interface for TribeNet reports starting with turn 899-12.
type tokenizer_899_12 struct {
	input  []byte
	length int
	pos    int
}

// getch returns the next byte in the input and advances the position.
// Returns 0 if at end of input.
func (tk *tokenizer_899_12) getch() byte {
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
func (tk *tokenizer_899_12) match(valid ...byte) (byte, bool) {
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
func (tk *tokenizer_899_12) next() Token {
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
	// note that 0x80, 0xD0, and 0xE2 are page separators in MS Word, so we treat them as whitespace
	if !glyphs[ch] {
		for tk.pos < tk.length && !glyphs[tk.peek()] {
			_ = tk.getch()
		}
		return Token{Type: Whitespace}
	}
	// single character tokens
	switch ch {
	case '&':
		return Token{Type: Ampersand}
	case '@':
		return Token{Type: AtSign}
	case '\\':
		return Token{Type: Backslash}
	case ':':
		return Token{Type: Colon}
	case ',':
		return Token{Type: Comma}
	case '-':
		return Token{Type: Dash}
	case '$':
		return Token{Type: DollarSign}
	case '.':
		return Token{Type: Dot}
	case '>':
		return Token{Type: GreaterThan}
	case '#':
		return Token{Type: Hash}
	case '(':
		return Token{Type: LeftParen}
	case ')':
		return Token{Type: RightParen}
	case ';':
		return Token{Type: Semicolon}
	case '/':
		return Token{Type: Slash}
	case '_':
		return Token{Type: Underscore}
	}
	// everything else is a text token
	if !delimiters[ch] {
		// if it's not a delimiter, keep reading until we hit one or whitespace
		for tk.pos < tk.length && glyphs[tk.peek()] && !delimiters[tk.peek()] {
			_ = tk.getch()
		}
	}
	// save the value as a string
	length := tk.pos - offset
	value := string(tk.input[offset : offset+length])
	// return it as a text token
	return Token{Type: Text, Value: value}
}

// peek returns the next byte in the input without advancing the position.
// Returns 0 if at end of input.
func (tk *tokenizer_899_12) peek() byte {
	if tk.pos >= tk.length {
		return 0
	}
	return tk.input[tk.pos]
}

// peekNext returns the next plus one byte in the input without advancing the position.
// Returns 0 if at end of input.
func (tk *tokenizer_899_12) peekNext() byte {
	if tk.pos+1 >= tk.length {
		return 0
	}
	return tk.input[tk.pos+1]
}
