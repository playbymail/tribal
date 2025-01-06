// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package pigeon

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }

const (
	ErrInvalidClan                     Error = "invalid clan"
	ErrInvalidFileType                 Error = "invalid file type"
	ErrInvalidReportFileName           Error = "invalid report file name"
	ErrInvalidTurnMonth                Error = "invalid turn month"
	ErrInvalidTurnNo                   Error = "invalid turn number"
	ErrInvalidTurnYear                 Error = "invalid turn year"
	ErrMissingTurnNumberInFirstSection Error = "missing turn number in first section"
	ErrMissingCurrentTurn              Error = "missing current turn"
	ErrNoData                          Error = "no data"
	ErrNoTurnNumber                    Error = "no turn number"
	ErrNoUnits                         Error = "no units found"
	ErrUnexpectedTurnNumber            Error = "unexpected turn number"
)

