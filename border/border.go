// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package border

import (
	"encoding/json"
	"fmt"
)

// Border_e is an enum for a border feature blocking an edge of a hex.
type Border_e int

const (
	None Border_e = iota
	Canal
	River
)

// MarshalJSON implements the json.Marshaler interface.
func (e Border_e) MarshalJSON() ([]byte, error) {
	return json.Marshal(EnumToString[e])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Border_e) UnmarshalJSON(data []byte) error {
	var s string
	var ok bool
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if *e, ok = StringToEnum[s]; !ok {
		return fmt.Errorf("invalid Border %q", s)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e Border_e) String() string {
	if str, ok := EnumToString[e]; ok {
		return str
	}
	return fmt.Sprintf("Border(%d)", int(e))
}

var (
	// EnumToString is a helper map for marshalling the enum
	EnumToString = map[Border_e]string{
		None:  "",
		Canal: "Canal",
		River: "River",
	}
	// StringToEnum is a helper map for unmarshalling the enum
	StringToEnum = map[string]Border_e{
		"":      None,
		"Canal": Canal,
		"River": River,
	}

	// LowerCaseToEnum is a helper map for parsing the border feature
	LowerCaseToEnum = map[string]Border_e{
		"canal": Canal,
		"river": River,
	}
)
