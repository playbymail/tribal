// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package stdlib

func IsTurn(year, month int) bool {
	switch {
	case year < 899:
		return false
	case year == 899:
		return month == 12
	case year < 9999:
		return 1 <= month && month <= 12
	default:
		return false
	}
}

func TurnFromYearMonth(year, month int) (int, bool) {
	if !IsTurn(year, month) {
		return 0, false
	}
	return ((year - 899) * 12) + month - 12, true
}

func TurnToYearMonth(turn int) (year, month int, ok bool) {
	if !(0 <= turn && turn < 10_000) {
		return 0, 0, false
	}

	// add 12 to the turn since we subtracted 12 in the original formula
	adjustedTurn := turn + 12

	// calculate year: divide by 12 and add back the base year (899)
	year = (adjustedTurn / 12) + 899

	// calculate month: get the remainder after dividing by 12
	// if remainder is 0, month should be 12
	month = adjustedTurn % 12
	if month == 0 {
		year, month = year-1, 12
	}

	return year, month, true
}
