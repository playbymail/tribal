// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package scanner

import (
	"fmt"
	"strings"
)

// patterns applies patches to the token stream to detect and correct common report errors.

// ApplyPatterns applies patterns to the token stream.
// bug: "Tribe 0987, , Current Hex" should keep those two commas.
func ApplyPatterns(head *Token) *Token {
	curr := head
	if curr == nil {
		panic("assert(curr != nil)")
	}
	for curr.Type != EOF {
		// replace whitespace+ with single whitespace
		for curr.next.Type == Whitespace && curr.next.next.Type == Whitespace {
			curr.next = curr.next.next
			curr.next.prev = curr
		}
		// replace comma (comma | spaces)+ with single comma.
		for curr.Type == Comma && (curr.next.Type == Comma || curr.next.Type == Whitespace) {
			curr.next = curr.next.next
			curr.next.prev = curr
		}
		// replace backslash (backslash | comma | dash | spaces)+ with single backslash
		for curr.Type == Backslash && (curr.next.Type == Comma || curr.next.Type == Dash || curr.next.Type == Whitespace) {
			curr.next = curr.next.next
			curr.next.prev = curr
		}
		// replace (comma | dash | whitespace) backslash with single backslash
		if curr.Type == Backslash && (curr.prev.Type == Comma || curr.prev.Type == Dash || curr.prev.Type == Whitespace) {
			if curr.prev.prev == nil {
				panic("assert(curr.prev.prev != nil)")
			}
			// remove the previous token
			curr.prev = curr.prev.prev
			curr.prev.next = curr
			// backtrack to the previous token so we can check other patterns that might have been affected
			curr = curr.prev
			continue
		}
		// replace (backslash | comma | spaces) newline with single newline
		if curr.Type == Newline && (curr.prev.Type == Backslash || curr.prev.Type == Comma || curr.prev.Type == Whitespace) {
			if curr.prev.prev == nil {
				panic("assert(curr.prev.prev != nil)")
			}
			// remove the previous token
			curr.prev = curr.prev.prev
			curr.prev.next = curr
			// backtrack to the previous token so we can check other patterns that might have been affected
			curr = curr.prev
			continue
		}

		// splice tokens together to create larger keywords and fix some issues for the parser.
		// it makes the parser easier to write, but the trade-off is that we can break settlement names.
		if MatchSequence(curr, Hash, Hash, Whitespace, Number) {
			hash := curr
			secondHash := hash.next
			whitespace := secondHash.next
			number := whitespace.next
			location := &Token{Type: Location, Value: fmt.Sprintf("## %s", number.String())}
			curr.prev.next = location
			location.next = number.next
			location.prev = hash.prev
			curr = location
		} else if MatchSequence(curr, Direction, Slash, Text) && strings.EqualFold(curr.Value, "N") && strings.EqualFold(curr.next.next.Value, "A") {
			curr.Type = Location
			curr.Value = "N/A"
			curr.next = curr.next.next.next
			curr.next.prev = curr
		} else if MatchSequence(curr, Tribe, Whitespace, Number, Comma, Whitespace, Comma) {
			// we have to change the type of the field for the tribe name to text so the comma compression doesn't delete it
			tribe := curr
			whitespace := tribe.next
			number := whitespace.next
			comma1 := number.next
			tribeName := comma1.next
			tribeName.Type = Text
			tribeName.Value = ""
		} else if MatchSequence(curr, Tribe, Whitespace, Number, Comma, Comma) {
			// we have to insert a field for the tribe name
			tribe := curr
			whitespace := tribe.next
			number := whitespace.next
			commaPre := number.next
			commaPost := commaPre.next
			tribeName := &Token{Type: Text, prev: commaPre, next: commaPost}
			commaPre.next = tribeName
			commaPost.prev = tribeName
		} else if MatchSequence(curr, Current, Whitespace, Hex) {
			curr.Type, curr.next = CurrentHex, curr.next.next.next
			curr.next.prev = curr
		} else if MatchSequence(curr, Goes, Whitespace, To) {
			curr.Type, curr.next = GoesTo, curr.next.next.next
			curr.next.prev = curr
		} else if MatchSequence(curr, Grassy, Whitespace, Hills) {
			curr.Type, curr.Value, curr.next = Terrain, "Grassy Hills", curr.next.next.next
			curr.next.prev = curr
		} else if MatchSequence(curr, Iron, Whitespace, Ore) {
			curr.Type, curr.Value, curr.next = Resource, "Iron Ore", curr.next.next.next
			curr.next.prev = curr
		} else if MatchSequence(curr, Jungle, Whitespace, Hills) {
			curr.Type, curr.Value, curr.next = Terrain, "Jungle Hills", curr.next.next.next
			curr.next.prev = curr
		} else if MatchSequence(curr, Low, Whitespace, Jungle, Whitespace, Mountains) {
			curr.Type, curr.Value, curr.next = Terrain, "Low Jungle Mountains", curr.next.next.next.next.next
			curr.next.prev = curr
		} else if MatchSequence(curr, Previous, Whitespace, Hex) {
			curr.Type, curr.next = PreviousHex, curr.next.next.next
		} else if MatchSequence(curr, Stone, Whitespace, Road) {
			curr.Type, curr.Value, curr.next = Terrain, "Stone Road", curr.next.next.next
			curr.next.prev = curr
		}

		// advance to the next token
		curr = curr.next
		if curr == nil {
			panic("assert(curr != nil)")
		}
	}
	return head
}

// MatchSequence returns true if the next tokens matches the sequence of given types.
func MatchSequence(curr *Token, types ...Type) bool {
	for len(types) != 0 && curr.Type == types[0] {
		types, curr = types[1:], curr.next
	}
	return len(types) == 0
}
