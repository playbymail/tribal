// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package lexer

var (
	keywords = map[string]Kind{
		"courier":  Courier,
		"current":  Current,
		"element":  Element,
		"fall":     Season,
		"fine":     Weather,
		"fleet":    Fleet,
		"hex":      Hex,
		"previous": Previous,
		"spring":   Season,
		"summer":   Season,
		"tribe":    Tribe,
		"turn":     Turn,
		"winter":   Season,
	}
)
