// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package norm

import "bytes"

// NormalizeCase returns a copy of the input with all characters converted to lower case.
func NormalizeCase(input []byte) []byte {
	return bytes.ToLower(input)
}
