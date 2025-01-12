// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package passage

import (
	"encoding/json"
	"fmt"
)

// Passage_e is an enum for a passage exiting a hex.
type Passage_e int

const (
	None Passage_e = iota
	Canal
	Ford
	Pass
	River
	StoneRoad
)

// MarshalJSON implements the json.Marshaler interface.
func (e Passage_e) MarshalJSON() ([]byte, error) {
	return json.Marshal(EnumToString[e])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Passage_e) UnmarshalJSON(data []byte) error {
	var s string
	var ok bool
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if *e, ok = StringToEnum[s]; !ok {
		return fmt.Errorf("invalid Edge %q", s)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e Passage_e) String() string {
	if str, ok := EnumToString[e]; ok {
		return str
	}
	return fmt.Sprintf("Edge(%d)", int(e))
}

var (
	// EnumToString is a helper map for marshalling the enum
	EnumToString = map[Passage_e]string{
		None:      "",
		Canal:     "Canal",
		Ford:      "Ford",
		Pass:      "Pass",
		River:     "River",
		StoneRoad: "Stone Road",
	}
	// StringToEnum is a helper map for unmarshalling the enum
	StringToEnum = map[string]Passage_e{
		"":           None,
		"Canal":      Canal,
		"Ford":       Ford,
		"Pass":       Pass,
		"River":      River,
		"Stone Road": StoneRoad,
	}

	// LowerCaseToEnum is a helper map for parsing the passage
	LowerCaseToEnum = map[string]Passage_e{
		"":           None,
		"canal":      Canal,
		"ford":       Ford,
		"pass":       Pass,
		"river":      River,
		"stone road": StoneRoad,
	}
)
