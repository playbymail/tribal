// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package adapters

import "github.com/playbymail/tribal"

func IntToTurnId(i int) (tribal.TurnId_t, bool) {
	if !(0 <= i && i < 9999) {
		return tribal.TurnId_t(0), false
	}
	return tribal.TurnId_t(i), true
}

// TurnIdToYearMonth converts a turn id to a year and month.
// It returns false if the turn id is invalid.
func TurnIdToYearMonth(turn tribal.TurnId_t) (year, month int, ok bool) {
	year, month = turn.YearMonth()
	return year, month, true
}

// YearMonthToTurnId converts a year and month to a turn id.
// It returns false if the year or month is invalid.
// Note that the year 899 is special and contains on the 12th month.
func YearMonthToTurnId(year, month int) (tribal.TurnId_t, bool) {
	if !(899 <= year && year <= 9999 && 1 <= month && month <= 12) {
		return 0, false
	}
	if year == 899 {
		return 0, month == 12
	}
	return tribal.TurnId_t((year-899)*12 + month - 12), true
}
