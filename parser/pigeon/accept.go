// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package pigeon

import (
	"fmt"
	"github.com/playbymail/tribal/parser/ast"
	"github.com/playbymail/tribal/parser/pigeon/rpt89912"
)

// acceptUnitId returns the unit id and true if the line starts with
// a unit type followed by a unit id and a comma.
func acceptUnitId(input []byte) (ast.Unit_t, bool) {
	if va, err := rpt89912.Parse("unit_id", input, rpt89912.Entrypoint("UnitId")); err != nil {
		// ignore the error. we know this means that a declaration was not found.
		return ast.Unit_t{}, false
	} else if unit, ok := va.(ast.Unit_t); !ok {
		panic(fmt.Sprintf("assert(%T == Unit_t)", va))
	} else {
		return unit, true
	}
}
