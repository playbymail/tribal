// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package goes_to

import (
	"github.com/playbymail/tribal/adapters"
	"github.com/playbymail/tribal/domains"
	"regexp"
)

var (
	// reUnitGoesTo is the regular expression for a unit goes to line.
	reUnitGoesTo = regexp.MustCompile(`^tribe goes to ([a-z]{2} \d{4})$`)
)

// Parse parses the tribe goes to line.
func Parse(path string, input []byte) (*domains.Coordinates_t, error) {
	if match := reUnitGoesTo.FindSubmatch(input); match == nil {
		return nil, domains.ErrInvalidUnitGoesTo
	} else if coords, err := adapters.TextToCoordinates(match[1]); err != nil {
		return nil, err
	} else {
		return &coords, err
	}
}
