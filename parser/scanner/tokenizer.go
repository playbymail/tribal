// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

// Tokenizer is the interface for tokenizing input.
// It should be implemented by any tokenizer and specialized to the version of the TribeNet report.
type Tokenizer interface {
	getch() byte
	match(...byte) (byte, bool)
	next() Token
	peek() byte
	peekNext() byte
}
