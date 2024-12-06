// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package turns

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }

const (
	ErrInvalidMonth  Error = "invalid month"
	ErrInvalidYear   Error = "invalid year"
	ErrInvalidTurnNo Error = "invalid turn number"
	ErrNoMatchFound  Error = "no match found"
)
