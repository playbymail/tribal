// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package turns_test

import (
	"github.com/playbymail/tribal/parser/turns"
	"testing"
)

// implements tests for parsing a turn line

func TestTurnLineParser(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected *turns.TurnLine_t
	}{
		{
			name:     "899-12",
			input:    `current turn 899-12(#0),winter,fine next turn 900-01(#1),29/10/2023`,
			expected: &turns.TurnLine_t{No: 0, Year: 899, Month: 12},
		},
		{
			name:     "900-05",
			input:    `current turn 900-05(#5),summer,fine next turn 900-06(#6),14/01/2024`,
			expected: &turns.TurnLine_t{No: 5, Year: 900, Month: 5},
		},
		{
			name:     "current only",
			input:    `current turn 900-05(#5),summer,fine`,
			expected: &turns.TurnLine_t{No: 5, Year: 900, Month: 5},
		},
		{
			name:  "invalid year",
			input: `current turn 876-05(#5),summer,fine`,
		},
	} {
		rslt, err := turns.ParseTurnLine(tc.name, []byte(tc.input))
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
			continue
		} else if rslt == nil && tc.expected == nil {
			continue
		} else if rslt == nil {
			t.Errorf("%s: expected result, got nil", tc.name)
			continue
		} else if tc.expected == nil {
			t.Errorf("%s: expected nil, got %+v", tc.name, *rslt)
			continue
		}

		if rslt.No != tc.expected.No {
			t.Errorf("%s: turn_no: expected %d, got %d", tc.name, tc.expected.No, rslt.No)
		}
		if rslt.Year != tc.expected.Year {
			t.Errorf("%s: year: expected %d, got %d", tc.name, tc.expected.Year, rslt.Year)
		}
		if rslt.Month != tc.expected.Month {
			t.Errorf("%s: month: expected %d, got %d", tc.name, tc.expected.Month, rslt.Month)
		}
	}
}
