// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

// Type is the type of token.
type Type int

const (
	Unknown Type = iota
	Ampersand
	AtSign
	Backslash
	BOF
	Colon
	Comma
	Dash
	DollarSign
	Dot
	EOF
	GreaterThan
	Hash
	HashHash
	LeftParen
	NA
	Newline
	RightParen
	Semicolon
	Slash
	Text
	Underscore
	Whitespace
)

type Token struct {
	Type  Type
	Value string
	// implement a linked list for the tokens
	prev *Token
	next *Token
}

func (t *Token) Prev() *Token {
	if t.Type == BOF {
		return t
	}
	return t.prev
}
func (t *Token) Next() *Token {
	if t.Type == EOF {
		return t
	}
	return t.next
}

func (t *Token) String() string {
	switch t.Type {
	case Ampersand:
		return "&"
	case AtSign:
		return "@"
	case Backslash:
		return "\\"
	case BOF:
		return ""
	case Colon:
		return ":"
	case Comma:
		return ","
	case Dash:
		return "-"
	case DollarSign:
		return "$"
	case Dot:
		return "."
	case EOF:
		return ""
	case GreaterThan:
		return ">"
	case Hash:
		return "#"
	case HashHash:
		return "##"
	case LeftParen:
		return "("
	case NA:
		return "N/A"
	case Newline:
		return "\n"
	case RightParen:
		return ")"
	case Semicolon:
		return ";"
	case Slash:
		return "/"
	case Text:
		return t.Value
	case Underscore:
		return "_"
	case Unknown:
		return "?"
	case Whitespace:
		return " "
	}
	return "?unknown?"
}
