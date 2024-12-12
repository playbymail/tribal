// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

import (
	"bytes"
)

// New returns a new scanner for the input.
// The list of tokens we buffer never includes EOF.
func New(input []byte) *Scanner {
	s := &Scanner{
		input:  input,
		pos:    0,
		length: len(input),
	}
	token := s.next()
	for ; token.Type != EOF; token = s.next() {
		s.tokens = append(s.tokens, token)
	}
	s.pos, s.input, s.length = 0, nil, len(s.tokens)
	return s
}

// Scanner tokenizes input bytes into a sequence of tokens.
// It buffers all tokens at creation time for efficient access.
// The tokens list never includes EOF.
type Scanner struct {
	input  []byte  // raw bytes being scanned (temporary, released after scanning)
	pos    int     // current position in input/tokens
	length int     // length of input/tokens
	tokens []Token // collected tokens from input
}

// Accept checks if the next token matches any of the given types.
// If there's a match, advances past that token and returns it with true.
// If no match, position remains unchanged and returns false.
// Returns false if at end of input.
func (s *Scanner) Accept(types ...Type) (Token, bool) {
	if s.pos >= s.length {
		return Token{}, false
	}
	t := s.tokens[s.pos]
	for _, v := range types {
		if t.Type == v {
			s.pos++
			return t, true
		}
	}
	return Token{}, false
}

// Backup steps back one token in the input.
// Will not back up past the beginning of input.
func (s *Scanner) Backup() {
	if s.pos > 0 {
		s.pos--
	}
}

// Next returns the next token in the input and advances the position.
// Returns EOF token if at end of input.
func (s *Scanner) Next() Token {
	if s.pos >= s.length {
		return Token{Type: EOF, offset: s.length}
	}
	token := s.tokens[s.pos]
	s.pos++
	return token
}

// Peek returns the next token in the input without advancing the position.
// Returns EOF token if at end of input.
func (s *Scanner) Peek() Token {
	if s.pos >= s.length {
		return Token{Type: EOF, offset: s.length}
	}
	return s.tokens[s.pos]
}

// PeekNext returns the next-next token in the input without advancing the position.
// Returns EOF token if there are fewer than 2 tokens remaining.
func (s *Scanner) PeekNext() Token {
	if s.pos+1 >= s.length {
		return Token{Type: EOF, offset: s.length}
	}
	return s.tokens[s.pos+1]
}

// RunOf returns a list of consecutive tokens that match any of the given types.
// Continues collecting tokens until it encounters one that doesn't match.
// Returns the collected sequence of matching tokens.
// Example: given input "a a b c" and type "a", returns ["a", "a"]
func (s *Scanner) RunOf(types ...Type) []Token {
	var run []Token
	for t := s.Peek(); t.Type != EOF; t = s.Peek() {
		isValid := false
		for _, v := range types {
			if v == t.Type {
				isValid = true
				break
			}
		}
		if !isValid {
			break
		}
		run = append(run, t)
		_ = s.Next()
	}
	return run
}

// RunTo returns a list of tokens up until one matches the given type.
// It does not include the token that matches the given type.
// At worst, it returns the entire input minus the EOF token.
// Example: given input "a b c d" and type "c", returns ["a", "b"]
func (s *Scanner) RunTo(types ...Type) []Token {
	var run []Token
	for t := s.Peek(); t.Type != EOF; t = s.Peek() {
		for _, to := range types {
			if to == t.Type {
				return run
			}
		}
		run = append(run, t)
		_ = s.Next()
	}
	return run
}

// Skip advances past the next token if it matches the given type.
// Returns true if a token was skipped, false if no match was found.
func (s *Scanner) Skip(skip Type) bool {
	_, ok := s.Accept(skip)
	return ok
}

// SkipRunOf skips a sequence of consecutive tokens that match any of the given types.
// Returns the number of tokens skipped.
func (s *Scanner) SkipRunOf(types ...Type) int {
	return len(s.RunOf(types...))
}

// SkipRunTo skips all tokens until finding one that matches any of the given types.
// Returns the number of tokens skipped, excluding the matching token.
func (s *Scanner) SkipRunTo(types ...Type) int {
	return len(s.RunTo(types...))
}

// Tokens returns all tokens collected from the input.
// The returned slice does not include the EOF token.
func (s *Scanner) Tokens() []Token {
	return s.tokens
}

// getch returns the next byte in the input and advances the position.
// Returns 0 if at end of input.
func (s *Scanner) getch() byte {
	if s.pos >= s.length {
		return 0
	}
	ch := s.input[s.pos]
	s.pos++
	return ch
}

// match checks if the next byte matches any in the valid set.
// If there's a match, advances position and returns the byte with true.
// If no match or at end of input, returns 0 with false.
func (s *Scanner) match(valid ...byte) (byte, bool) {
	if s.pos >= s.length {
		return 0, false
	}
	ch := s.input[s.pos]
	if bytes.IndexByte(valid, ch) == -1 {
		return 0, false
	}
	s.pos++
	return ch, true
}

// peek returns the next byte in the input without advancing the position.
// Returns 0 if at end of input.
func (s *Scanner) peek() byte {
	if s.pos >= s.length {
		return 0
	}
	return s.input[s.pos]
}

// next returns the next token in the input.
// Handles newlines, whitespace, and text tokens distinctly.
// For text tokens, converts value to lowercase.
// Returns EOF token if at end of input.
func (s *Scanner) next() Token {
	if s.pos >= s.length {
		return Token{Type: EOF, offset: s.length}
	}
	// anchor the token
	token := Token{offset: s.pos}
	ch := s.getch()
	if ch == '\n' {
		token.Type = Newline
		token.length = s.pos - token.offset
		return token
	}
	// lump whitespace and invalid characters together
	if !glyphs[ch] {
		token.Type = Whitespace
		for s.pos < s.length && !glyphs[s.peek()] {
			_ = s.getch()
		}
		token.length = s.pos - token.offset
		return token
	}
	token.Type = Text
	if !delimiters[ch] {
		for s.pos < s.length && glyphs[s.peek()] && !delimiters[s.peek()] {
			_ = s.getch()
		}
	}
	token.length = s.pos - token.offset
	token.Value = string(bytes.ToLower(s.input[token.offset : token.offset+token.length]))
	return token
}
