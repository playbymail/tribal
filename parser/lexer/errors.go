// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package lexer

const (
	ErrNoInput Error = "no input"
)

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }
