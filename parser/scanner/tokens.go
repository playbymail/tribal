// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

import "regexp"

// there are a few tokens that are context-sensitive.
// for example, SW can mean "South West" or "Swamp."
// the digits 1214 can be a number, a tribe id, or a hex location.
// we can't guess, so we choose the most likely interpretation
// based on usage in the report. we trust the parser to
// use its context to disambiguate.

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
	Courier
	CourierID
	Current
	CurrentHex
	Dash
	Direction
	DollarSign
	Dot
	DoubleQuote
	Element
	ElementID
	EOF
	Fleet
	FleetID
	Follows
	Garrison
	GarrisonID
	Goes
	GoesTo
	Grassy
	GreaterThan
	GridID
	Hash
	HashHash
	Hex
	Hills
	Iron
	Jungle
	Land
	LeftParen
	Location
	Low
	Mountains
	Move
	Movement
	Newline
	Number
	Ore
	Previous
	PreviousHex
	Resource
	RightParen
	Road
	Scout
	Semicolon
	Sight
	SingleQuote
	Slash
	Stone
	Terrain
	Text
	To
	Tribe
	TribeID
	Turn
	Underscore
	Water
	Whitespace
	Winds
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
	case Courier:
		return "Courier"
	case CourierID:
		return t.Value
	case Current:
		return "Current"
	case CurrentHex:
		return "Current Hex"
	case Dash:
		return "-"
	case Direction:
		return t.Value
	case DollarSign:
		return "$"
	case Dot:
		return "."
	case DoubleQuote:
		return "\""
	case Element:
		return "Element"
	case ElementID:
		return t.Value
	case EOF:
		return ""
	case Fleet:
		return "Fleet"
	case FleetID:
		return t.Value
	case Follows:
		return "follows"
	case Garrison:
		return "Garrison"
	case GarrisonID:
		return t.Value
	case Goes:
		return "goes"
	case GoesTo:
		return "goes to"
	case Grassy:
		return "Grassy"
	case GreaterThan:
		return ">"
	case GridID:
		return t.Value
	case Hash:
		return "#"
	case HashHash:
		return "##"
	case Hex:
		return "Hex"
	case Hills:
		return "Hills"
	case Iron:
		return "Iron"
	case Jungle:
		return "Jungle"
	case Land:
		return "Land"
	case LeftParen:
		return "("
	case Location:
		return t.Value
	case Low:
		return "Low"
	case Mountains:
		return "Mountains"
	case Move:
		return "Move"
	case Movement:
		return "Movement"
	case Newline:
		return "\n"
	case Number:
		return t.Value
	case Ore:
		return "Ore"
	case Previous:
		return "Previous"
	case PreviousHex:
		return "Previous Hex"
	case Resource:
		return t.Value
	case RightParen:
		return ")"
	case Road:
		return "Road"
	case Scout:
		return "Scout"
	case Semicolon:
		return ";"
	case Sight:
		return "Sight"
	case SingleQuote:
		return "'"
	case Slash:
		return "/"
	case Stone:
		return "Stone"
	case Terrain:
		return t.Value
	case Text:
		return t.Value
	case To:
		return "to"
	case Tribe:
		return "Tribe"
	case TribeID:
		return t.Value
	case Turn:
		return "Turn"
	case Underscore:
		return "_"
	case Unknown:
		return "?"
	case Water:
		return "Water"
	case Whitespace:
		return " "
	case Winds:
		return t.Value
	}
	return "?unknown?"
}

var (
	keywords = map[string]Type{
		"canal":     Terrain,
		"courier":   Courier,
		"current":   Current,
		"element":   Element,
		"fleet":     Fleet,
		"follows":   Follows,
		"ford":      Terrain,
		"garrison":  Garrison,
		"goes":      Goes,
		"grassy":    Grassy,
		"hex":       Hex,
		"hills":     Hills,
		"iron":      Iron,
		"jungle":    Jungle,
		"l":         Terrain,
		"lake":      Terrain,
		"low":       Low,
		"mild":      Winds,
		"mountains": Mountains,
		"move":      Move,
		"movement":  Movement,
		"n":         Direction,
		"ne":        Direction,
		"nw":        Direction,
		"o":         Terrain,
		"ocean":     Terrain,
		"ore":       Ore,
		"prairie":   Terrain,
		"previous":  Previous,
		"river":     Terrain,
		"road":      Road,
		"s":         Direction,
		"scout":     Scout,
		"se":        Direction,
		"sight":     Sight,
		"stone":     Stone,
		"sw":        Direction,
		"to":        To,
		"tribe":     Tribe,
		"turn":      Turn,
	}

	rxCourierID  = regexp.MustCompile(`^\d{4}c[1-9]$`)
	rxElementID  = regexp.MustCompile(`^\d{4}e[1-9]$`)
	rxFleetID    = regexp.MustCompile(`^\d{4}f[1-9]$`)
	rxGarrisonID = regexp.MustCompile(`^\d{4}g[1-9]$`)
	rxNumber     = regexp.MustCompile(`^\d+$`)
	rxTribeID    = regexp.MustCompile(`^\d{4}$`)
)
