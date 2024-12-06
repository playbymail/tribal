// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package adapters

import (
	"github.com/playbymail/tribal"
)

func IntToClanId(i int) (tribal.ClanId_t, bool) {
	if !(0 < i && i < 1000) {
		return tribal.ClanId_t(0), false
	}
	return tribal.ClanId_t(i), true
}
