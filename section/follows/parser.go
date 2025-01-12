// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package follows

import (
	"github.com/playbymail/tribal/domains"
	"regexp"
)

var (
	// reUnitFollows is the regular expression for a unit follows line.
	reUnitFollows = regexp.MustCompile(`^tribe follows (\d{4}(?:[cefg][1-9])?)$`)
)

// Parse parses the tribe follows line.
func Parse(path string, input []byte) (domains.UnitId_t, error) {
	match := reUnitFollows.FindSubmatch(input)
	if match == nil {
		return "", domains.ErrInvalidUnitFollows
	}
	return domains.UnitId_t(match[1]), nil
}
