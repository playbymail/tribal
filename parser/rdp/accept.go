// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package rdp

import (
	"bytes"
	"regexp"
)

var (
	rxCourierHeader  = regexp.MustCompile(`^courier (\d{4}c\d),`)
	rxElementHeader  = regexp.MustCompile(`^element (\d{4}e\d),`)
	rxFleetHeader    = regexp.MustCompile(`^fleet (\d{4}f\d),`)
	rxGarrisonHeader = regexp.MustCompile(`^garrison (\d{4}g\d),`)
	rxTribeHeader    = regexp.MustCompile(`^tribe (\d{4}),`)
)

// AcceptComma reads a comma from the input.
func AcceptComma(b []byte) (token, rest []byte) {
	if len(b) == 0 || b[0] != ',' {
		return nil, nil
	}
	return b[:1], b[1:]
}

// AcceptJunk reads to the next section header.
func AcceptJunk(b []byte) (token, rest []byte) {
	for i := 0; i < len(b); i++ {
		// check for a section header.
		// if we find one, return the bytes up to it, but not including it.
		if rxCourierHeader.Match(b[i:]) {
			return b[:i], b[i:]
		} else if rxElementHeader.Match(b[i:]) {
			return b[:i], b[i:]
		} else if rxFleetHeader.Match(b[i:]) {
			return b[:i], b[i:]
		} else if rxGarrisonHeader.Match(b[i:]) {
			return b[:i], b[i:]
		} else if rxTribeHeader.Match(b[i:]) {
			return b[:i], b[i:]
		}
		// no section header here, so skip to the end of this line.
		for i < len(b) && b[i] != LF {
			i++
		}
		// the LF (if any) will be consumed by the outer loop.
	}
	// remainder of buffer is all junk
	return b, nil
}

// AcceptToLF reads to the next LF or end of input.
// never includes the LF as part of the token.
// if there is no LF in the buffer, returns nil and the buffer.
// if the first byte of the buffer is a LF, returns an empty slice and the buffer.
// otherwise, returns a slice of bytes up to the LF
// and a slice containing the LF and remaining bytes.
func AcceptToLF(b []byte) (token, rest []byte) {
	offset := bytes.IndexByte(b, LF)
	if offset == -1 { // no LF in the buffer
		return nil, b
	}
	return b[:offset], b[offset:]
}
