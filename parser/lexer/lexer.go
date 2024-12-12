// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package lexer

import (
	"bytes"
	"regexp"
)

type Lexer struct {
	input    []byte // original input, caller owns this memory
	pos      int    // current position in input
	lexemes  []Lexeme
	pushback []Lexeme
}

type Lexeme struct {
	Kind   Kind
	Offset int // offset of lexeme in input
	Length int // length of lexeme in input
}
type Kind int

// New returns a new Lexer.
// The caller owns the input memory and must ensure it is valid for the lifetime of the Lexer.
// It is an error to pass an empty or nil input to the Lexer.
func New(input []byte) (*Lexer, error) {
	if len(input) == 0 {
		return nil, ErrNoInput
	}
	l := &Lexer{
		input:   input,
		lexemes: make([]Lexeme, 0, 600_000),
	}
	// for some magic to happen, we have to push a newline into the lexeme stream
	l.lexemes = append(l.lexemes, Lexeme{Kind: EOL, Offset: 0, Length: 0})
	for lexeme := l.next(); lexeme.Kind != EOF; lexeme = l.next() {
		l.lexemes = append(l.lexemes, lexeme)
	}
	return l, nil
}

func (l *Lexer) Next() Lexeme {
	if len(l.lexemes) == 0 {
		return Lexeme{Kind: EOF, Offset: len(l.input), Length: 0}
	}
	lexeme := l.lexemes[0]
	l.lexemes = l.lexemes[1:]
	return lexeme
}

func (l *Lexer) next() Lexeme {
	// skip whitespace
	for l.pos < len(l.input) && whitespace[l.input[l.pos]] {
		l.pos++
	}

	// if we're at the end of the input, return EOF
	if !(l.pos < len(l.input)) {
		return Lexeme{Kind: EOF, Offset: len(l.input), Length: 0}
	}

	// anchor the start of the lexeme
	lexeme := Lexeme{Offset: l.pos}
	start := l.pos

	// fetch the next character
	ch := l.getch()

	// return a run of invalid glyphs as a single lexeme
	if !validGlyphs[ch] {
		for ch = l.peek(); ch != 0 && !validGlyphs[ch]; ch = l.peek() {
			l.getch()
		}
		lexeme.Kind, lexeme.Length = InvalidGlyphs, l.pos-start
		return lexeme
	}

	// do we have a single-character delimiter?
	if d := singleCharDelimiters[ch]; d != EOF {
		lexeme.Kind, lexeme.Length = d, l.pos-start
		return lexeme
	} else if ch == '#' {
		// delimiter can be '##' or just '#'
		if _, ok := l.match('#'); ok {
			lexeme.Kind, lexeme.Length = HashHash, l.pos-start
			return lexeme
		}
		lexeme.Kind, lexeme.Length = Hash, l.pos-start
		return lexeme
	}

	//// some keyword checks
	//if ch == 'n' && l.peek() == '/' && l.peekNext() == 'a' {
	//	l.getch()
	//	l.getch()
	//	lexeme.Kind, lexeme.Length = NA, l.pos-start
	//	return lexeme
	//}

	// no matches, so treat it as text
	lexeme.Kind, lexeme.Length = Text, l.pos-start
	for textGlyphs[l.peek()] {
		l.getch()
	}
	lexeme.Length = l.pos - start

	// check for keywords
	word := l.input[lexeme.Offset : lexeme.Offset+lexeme.Length]

	if rxInteger.Match(word) {
		lexeme.Kind = Integer
		return lexeme
	} else if rxReportDate.Match(word) {
		lexeme.Kind = ReportDate
		return lexeme
	} else if rxTurnId.Match(word) {
		lexeme.Kind = TurnId
		return lexeme
	} else if rxUnitId.Match(word) {
		lexeme.Kind = UnitId
		return lexeme
	}

	text := string(bytes.ToLower(l.input[lexeme.Offset : lexeme.Offset+lexeme.Length]))
	if kind, ok := keywords[text]; ok {
		lexeme.Kind = kind
		return lexeme
	}

	// maybe a unit id

	return lexeme

	//for pos := 0; pos < len(input); {
	//
	//	// skip invalid characters
	//	if !validGlyphs[ch] {
	//		pos++
	//		continue
	//	}
	//
	//	// skip whitespace between tokens
	//	if whitespace[ch] {
	//		pos++
	//		continue
	//	}
	//
	//	// do we have a delimiter?
	//	if d := delimiters[ch]; d != EOF {
	//		l.lexemes = append(l.lexemes, Lexeme{Kind: d})
	//		pos++
	//		continue
	//	}
	//
	//	// not a simple delimiter, so look for a keyword. collect characters up to eof or a delimiter.
	//	offset := pos
	//	for offset < len(input) && delimiters[input[offset]] == EOF {
	//		offset++
	//	}
	//
	//	// no matches, so treat it as a settlement; collect characters up to eof or a delimiter
	//	for offset < len(input) && !settlementDelimiters[input[offset]] {
	//		offset++
	//	}
	//	l.lexemes = append(l.lexemes, Lexeme{
	//		Kind:  UnitName,
	//		Value: input[pos:offset],
	//	})
	//	pos = offset
	//}
	//return l
}

// getch returns the next character in the input.
// forces ch to lowercase as we're reading.
// returns 0 if we're at the end of the input.
func (l *Lexer) getch() (ch byte) {
	if !(0 <= l.pos && l.pos < len(l.input)) {
		return 0
	}
	ch, l.pos = l.input[l.pos], l.pos+1
	return tolower(ch)
}

func tolower(ch byte) byte {
	if 'A' <= ch && ch <= 'Z' {
		return ch - 'A' + 'a'
	}
	return ch
}

// match compares the next character in the input to the list of characters.
// returns the character and true if it is in the list.
func (l *Lexer) match(chars ...byte) (byte, bool) {
	if l.iseof() {
		return 0, false
	}
	nextChar := tolower(l.input[l.pos])
	for _, ch := range chars {
		if ch == nextChar {
			return l.getch(), true
		}
	}
	return 0, false
}

// iseof returns true if we're at the end of the input.
func (l *Lexer) iseof() bool {
	return !(l.pos < len(l.input))
}

// peek returns the next character in the input without advancing the position.
// returns 0 if we're at the end of the input or position is invalid.
func (l *Lexer) peek() (ch byte) {
	if !(0 <= l.pos && l.pos < len(l.input)) {
		return 0
	}
	return tolower(l.input[l.pos])
}

// peekNext returns the second character in the input without advancing the position.
// returns 0 if we're at the end of the input or position is invalid.
func (l *Lexer) peekNext() (ch byte) {
	if !(l.pos+1 < len(l.input)) {
		return 0
	}
	return tolower(l.input[l.pos+1])
}

// skipto returns the next occurrence of a character in the input.
// returns 0 we reach end of input before finding the character.
func (l *Lexer) skipto(ch byte) byte {
	// skip the current character.
	if 0 <= l.pos && l.pos < len(l.input) {
		l.pos++
	}
	for 0 <= l.pos && l.pos < len(l.input) && tolower(l.input[l.pos]) != ch {
		l.pos++
	}

	return ch
}

// ungetch returns the last character read back to the input.
// there's a certain wonkiness if we push the position to before,
// the start of the input, so we must ensure that the position is
// valid after updating. meaning 0 <= pos < len(input).
// we do assume that the input is never empty.
func (l *Lexer) ungetch() {
	l.pos = l.pos - 1
	if l.pos < 0 {
		l.pos = 0
	}
}

var (
	textGlyphs           [256]bool
	validGlyphs          [256]bool
	singleCharDelimiters [256]Kind
	settlementDelimiters [256]bool
	whitespace           [256]bool
)

func init() {
	// initialize the lookup table for whitespace
	for _, ch := range []byte(" \t") {
		whitespace[ch] = true
	}
	// initialize the lookup tables for glyphs
	for ch := '0'; ch <= '9'; ch++ {
		textGlyphs[ch], validGlyphs[ch] = true, true
	}
	for ch := 'A'; ch <= 'Z'; ch++ {
		textGlyphs[ch], validGlyphs[ch] = true, true
	}
	for ch := 'a'; ch <= 'z'; ch++ {
		textGlyphs[ch], validGlyphs[ch] = true, true
	}
	for _, ch := range []byte("-/$'.") {
		textGlyphs[ch], validGlyphs[ch] = true, true
	}
	for _, ch := range []byte("(){}[]<>`#%&*+,:;=?@^_|~\\\" \t\n") {
		validGlyphs[ch] = true
	}
	// initialize the lookup table for settlement delimiters
	for _, ch := range []byte(",\n") {
		settlementDelimiters[ch] = true
	}
	// initialize the lookup table for single-character delimiters
	singleCharDelimiters['\\'] = Backslash
	singleCharDelimiters[':'] = Colon
	singleCharDelimiters[','] = Comma
	singleCharDelimiters['-'] = Dash
	singleCharDelimiters['\n'] = EOL
	singleCharDelimiters['='] = Equals
	singleCharDelimiters['('], singleCharDelimiters[')'] = LParen, RParen
}

const (
	EOF                 Kind = iota
	Backslash                // Literal backslash
	BlockingTerrain          // Token for blocking terrain (e.g., "Lake", "Ocean", "River", "Swamp/Jungle Hill")
	CANT                     // Literal "Can't"
	CANT_MOVE_INTO           // Literal "Cannot Move Wagons into"
	CANT_MOVE_ON             // Literal "Can't Move on"
	CANNOT                   // Literal "Cannot"
	Colon                    // Literal ':'
	Comma                    // Literal ','
	Courier                  // Literal 'Courier'
	COURIER_UNIT_ID          // Token for courier unit
	Current                  // Literal "Current"
	CURRENT_HEX              // Literal 'Current Hex'
	Dash                     // Literal '-'
	DIRECTION                // Token for direction (e.g., "N", "SE")
	DIRECTION_CROWSNEST      // Token for direction from outer ring (e.g., "N/N", "N/NE")
	EDGE_CSV                 // Token for edge terrain with comma separated list (e.g. "Hsm", "L", "Lcm", "O")
	EDGE_LIST                // Token for edge terrain with space separated list (e.g. "Ford", "Pass", "River")
	EOL                      // Literal end of line
	Equals                   // Literal '='
	Element                  // Literal 'Element'
	ELEMENT_UNIT_ID          // Token for element unit
	EXHAUSTED                // Literal "Not enough M.P's to move to"
	FIND                     // Literal 'Find'
	Fleet                    // Literal 'Fleet'
	FLEET_UNIT_ID            // Token for fleet unit
	FOLLOWS                  // Literal 'Follows'
	Ford                     // Literal "Ford"
	Garrison                 // Literal 'Garrison'
	GARRISON_UNIT_ID         // Token for garrison unit
	GOES_TO                  // Literal 'Goes to'
	GRID_LOCATION            // Token for location on the grid map (e.g., "AA 1010")
	Hash                     // Literal '#'
	HashHash                 // Literal "##", which is usually an obscured grid id
	Hex                      // Literal "Hex"
	Integer                  // Token for integer
	InvalidGlyphs            // the token contains invalid glyphs
	INTO                     // Literal "into"
	ITEM                     // Token for items (e.g., "Frame", "Silver")
	LAND                     // Literal 'Land'
	LParen                   // Literal '('
	MOVE                     // Literal 'Move'
	MOVEMENT                 // Literal 'Movement'
	NA                       // Literal 'N/A'
	NO_FORD                  // Literal "No Ford on"
	NO_RIVER                 // Literal "No River Adjacent to Hex to"
	NOTHING_FOUND            // Literal "Nothing of interest found"
	Ocean                    // Literal "Ocean"
	OF_HEX                   // Literal 'of Hex'
	ON                       // Literal "on"
	PATROLLED_AND_FOUND      // Literal "Patrolled and found"
	Previous                 // Literal "Previous"
	PREVIOUS_HEX             // Literal 'Previous Hex'
	QUANTITY                 // Token for number of items (e.g., 1)
	ReportDate               // Token for report date (dd/mm/yyyy)
	RESOURCE                 // Token for resources (e.g., "IRON ORE")
	River                    // Literal "River"
	RParen                   // Literal ')'
	Terrain                  // Token for terrain (e.g., "PR", "SW")
	TerrainName              // Token for terrain name (e.g., "DECIDUOUS", "ROCKY HILLS", "SWAMP")
	Text                     // Token for text that doesn't match a keyword (e.g., "A", "B", "C")
	TO                       // Literal 'to'
	Tribe                    // Literal 'Tribe'
	TRIBE_UNIT_ID            // Token for tribe unit
	Turn                     // Literal "Turn"
	TurnId                   // Token for turn number (year-month)
	SCOUT                    // Literal 'Scout'
	ScoutId                  // Token for scout number (e.g., 1, 2)
	Season                   // Token for season (e.g., "Winter", "Summer")
	Settlement               // Token for settlement name (really, token for anything that isn't another token)
	SIGHT                    // Literal 'Sight'
	STATUS                   // Literal 'Status'
	STILL                    // Literal 'Still'
	UnitId                   // A token for the unit's unique identifier
	UnitName                 // A token for the unit's name
	WAGONS                   // Literal "Wagons"
	WATER                    // Literal 'Water'
	Weather                  // Token for weather (e.g., "FINE", "RAINY")
	WINDS                    // A token for wind strength (e.g., "MILD")
)

var (
	rxInteger    = regexp.MustCompile(`^\d+$`)
	rxReportDate = regexp.MustCompile(`^([1-9]|[1-9]\d)/(1[0-2]|[1-9])/\d{4}$`)
	rxTurnId     = regexp.MustCompile(`^\d\d\d-\d\d$`)
	rxUnitId     = regexp.MustCompile(`^\d{4}[cefg][1-9]$`)
)
