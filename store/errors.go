// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package store

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }

const (
	ErrDatabase              = Error("database error")
	ErrDuplicateClanId Error = "duplicate clan id"
	ErrDuplicateReport Error = "duplicate report"
	ErrInvalidClanId   Error = "invalid clan id"
	ErrInvalidMonth    Error = "invalid month"
	ErrInvalidTurnNo   Error = "invalid turn no"
	ErrInvalidYear     Error = "invalid year"
	ErrNoData          Error = "no data"
	ErrNotExist        Error = "database file does not exist"
	ErrNotFound        Error = "not found"
	ErrNotImplemented  Error = "not implemented"
)
