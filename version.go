// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/semver"
)

var (
	version = semver.Version{Major: 0, Minor: 1, Patch: 0}
)

func Version() semver.Version {
	return version
}
