// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package stdlib

import (
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	reReportFile = regexp.MustCompile(`^([0-9]{3,4})-([0-9]{1,2})\.([0-9]{3,4})\.report\.(docx|txt)$`)
)

// IsReportFileName verifies that the path contains a valid report file name.
// Because of Raven, we don't require an exact match. We accept the name if
// it is kind of close to the YYYY-MM.CLAN.report pattern. We accept only
// docx and txt for the extension.
// Returns false if the name is not valid.
// Otherwise, returns the clan, year, and month extracted from the name.
func IsReportFileName(path string) (clan, year, month int, ok bool) {
	//log.Printf("path: %s", path)
	name := filepath.Base(path)
	//log.Printf("name: %s", name)
	matches := reReportFile.FindStringSubmatch(name)
	//log.Printf("matches: %v", matches)
	if len(matches) == 0 {
		return 0, 0, 0, false
	}
	// extract and validate the year
	if n, err := strconv.Atoi(matches[1]); err != nil {
		//log.Printf("error parsing year: %q: %v", matches[1], err)
		return 0, 0, 0, false
	} else if !(899 <= n && n <= 9999) {
		//log.Printf("year out of range: %q: %d", matches[1], n)
		return 0, 0, 0, false
	} else {
		year = n
	}
	// extract and validate the month
	if n, err := strconv.Atoi(matches[2]); err != nil {
		//log.Printf("error parsing month: %q: %v", matches[2], err)
		return 0, 0, 0, false
	} else if !(1 <= n && n <= 12) {
		//log.Printf("month out of range: %q: %d", matches[2], n)
		return 0, 0, 0, false
	} else {
		month = n
	}
	// extract and validate the clan
	if n, err := strconv.Atoi(matches[3]); err != nil {
		//log.Printf("error parsing clan: %q: %v", matches[3], err)
		return 0, 0, 0, false
	} else if !(0 < n && n < 1000) {
		//log.Printf("clan out of range: %q: %d", matches[3], n)
		return 0, 0, 0, false
	} else {
		clan = n
	}
	return clan, year, month, true
}
