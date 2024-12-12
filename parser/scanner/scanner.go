// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package scanner implements a lexical scanner for TribeNet reports.
package scanner

import "fmt"

// New returns a new scanner for the report input.
// Version expects a string of the form "899-12" or "902-05."
func New(input []byte, version string) (*Scanner, error) {
	var tk Tokenizer
	switch version {
	case "":
		return nil, fmt.Errorf("missing version", version)
	case "899-12":
		tk = &tokenizer_899_12{input: input, length: len(input)}
	default:
		return nil, fmt.Errorf("%s: unsupported version", version)
	}
	// create a scanner from the tokenizer.
	s := &Scanner{
		head: &Token{Type: BOF},
	}
	// set current and tail to the head of the list
	s.curr, s.tail = s.head, s.head
	// add remaining tokens to the list
	for token := tk.next(); token.Type != EOF; token = tk.next() {
		token.prev = s.tail
		s.tail.next = &token
		s.tail = s.tail.next
		s.length++
	}
	// ensure that we have an EOF token at the end of the list
	s.tail.next = &Token{Type: EOF, prev: s.tail}
	s.tail = s.tail.next
	return s, nil
}

// Scanner tokenizes input bytes into a sequence of tokens.
// It buffers all tokens at creation time for efficient access.
type Scanner struct {
	head   *Token
	curr   *Token
	tail   *Token
	length int // number of tokens in the list
}

// Accept checks if the next token matches any of the given types.
// If there's a match, advances past that token and returns it with true.
// If no match, position remains unchanged and returns false.
// Returns false if at end of input.
func (s *Scanner) Accept(types ...Type) (*Token, bool) {
	if s.curr.Type == EOF {
		return nil, false
	}
	t := s.curr
	for _, v := range types {
		if t.Type == v {
			s.curr = s.curr.next
			return t, true
		}
	}
	return nil, false
}

// Backup steps back one token in the input.
// Will not back up past the beginning of input.
func (s *Scanner) Backup() {
	if s.curr.prev == nil {
		return
	}
	s.curr = s.curr.prev
}

// Next returns the next token in the input and advances the position.
// Returns EOF token if at end of input.
func (s *Scanner) Next() *Token {
	t := s.curr
	if s.curr.Type != EOF {
		s.curr = s.curr.next
	}
	return t
}

func (s *Scanner) NumTokens() int {
	return s.length
}

// Peek returns the next token in the input without advancing the position.
// Returns EOF token if at end of input.
func (s *Scanner) Peek() *Token {
	return s.curr
}

// PeekNext returns the next-next token in the input without advancing the position.
// Returns EOF token if there are fewer than 2 tokens remaining.
func (s *Scanner) PeekNext() *Token {
	if s.curr.Type == EOF {
		return s.curr
	}
	return s.curr.next
}

// RunOf returns a list of consecutive tokens that match any of the given types.
// Continues collecting tokens until it encounters one that doesn't match.
// Returns the collected sequence of matching tokens.
// Example: given input "a a b c" and type "a", returns ["a", "a"]
func (s *Scanner) RunOf(types ...Type) []*Token {
	var run []*Token
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
func (s *Scanner) RunTo(types ...Type) []*Token {
	var run []*Token
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
