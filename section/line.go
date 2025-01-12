// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package section

type Line struct {
	Line int      // line number in the original input
	Kind LineKind // follows, goes to, fleet, land, scouts, status
	Text []byte   // text of the move
}

// Less returns true if line should be sorted before the other line.
func (l *Line) Less(l2 *Line) bool {
	// sort by kind, then by line number
	if l.Kind < l2.Kind {
		return true
	} else if l.Kind == l2.Kind {
		return l.Line < l2.Line
	}
	return false
}

// warning - lines are sorted by kind first, then by line number,
// so keep that in mind when updating the LinkKind enum.

type LineKind int

const (
	Unknown LineKind = iota
	Turn
	UnitFollows
	UnitGoesTo
	UnitLandMove
	UnitFleetMove
	UnitScouts
	UnitStatus
)
