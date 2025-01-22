// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ast

import "encoding/json"

// HexName_t is a special hex name or a village name.
// We don't actually know how to distinguish between the two in the parser,
// so we default to a village name since it is more common. At some point,
// the front end will have to allow the user to override the default.
type HexName_t struct {
	Type HexName_e `json:"type"`
	Name string    `json:"name"`
}

type HexName_e int // enum for hex name type
const (
	VillageName HexName_e = iota
	SpecialHex
)

func (e HexName_e) MarshalJSON() ([]byte, error) {
	if e == VillageName {
		return json.Marshal("Village")
	}
	return json.Marshal("Special Hex")
}

func (e HexName_e) MarshalText() ([]byte, error) {
	if e == VillageName {
		return []byte("Village"), nil
	}
	return []byte("Special Hex"), nil
}
