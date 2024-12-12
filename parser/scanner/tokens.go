// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

type Token struct {
	Type   Type
	Value  string
	offset int // offset of token in input
	length int // length of token in input
}

func (t Type) String() string {
	switch t {
	case EOF:
		return "EOF"
	case Newline:
		return "CR"
	case Text:
		return "TXT"
	case Whitespace:
		return "SP"
	case Unknown:
		return "UNK"
	default:
		return "???"
	}
}

// Type is the type of token.
type Type int

const (
	Unknown Type = iota
	EOF
	Newline
	Text
	Whitespace
)
