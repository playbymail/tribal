// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package ast

const (
	ErrInvalidCoordinates    Error = "invalid coordinates"
	ErrMultipleCurrentHexes  Error = "multiple current hexes"
	ErrMultiplePreviousHexes Error = "multiple previous hexes"
	ErrNoMatch               Error = "no match"
	ErrTooFewFields          Error = "too few fields"
	ErrTooManyFields         Error = "too many fields"
	ErrUnexpectedInput       Error = "unexpected input"
)

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }
