// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

const (
	ErrInvalidCoordinates    Error = "invalid coordinates"
	ErrInvalidMonth          Error = "invalid month"
	ErrInvalidTurnNo         Error = "invalid turn number"
	ErrInvalidUnitFollows    Error = "invalid unit follows"
	ErrInvalidUnitGoesTo     Error = "invalid unit goes to"
	ErrInvalidYear           Error = "invalid year"
	ErrMissingTerrainType    Error = "missing terrain type"
	ErrMultipleCurrentHexes  Error = "multiple current hexes"
	ErrMultiplePreviousHexes Error = "multiple previous hexes"
	ErrNoMatch               Error = "no match"
	ErrNotTribeMovementLine  Error = "not a tribe movement line"
	ErrNotUnitStatusLine     Error = "not a unit status line"
	ErrTooFewFields          Error = "too few fields"
	ErrTooManyFields         Error = "too many fields"
	ErrTurnNoMismatch        Error = "turn number mismatch"
	ErrUnexpectedInput       Error = "unexpected input"
)

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }
