// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package common

import (
	"github.com/playbymail/tribal/parser/ast"
	"regexp"
)

var (
	// reUnitFollows is the regular expression for a unit follows line.
	reUnitFollows = regexp.MustCompile(`^tribe follows (\d{4}(?:[cefg][1-9])?)$`)
)

// ParseTribeFollows parses the tribe follows line.
//
//	"tribe follows" UnitId
func ParseTribeFollows(path string, input []byte) (ast.UnitId_t, error) {
	match := reUnitFollows.FindSubmatch(input)
	if match == nil {
		return "", ast.ErrInvalidUnitFollows
	}
	return ast.UnitId_t(match[1]), nil
}
