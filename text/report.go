// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package text

import (
	"github.com/playbymail/tribal/norm"
	"os"
	"strings"
)

// ReadFile loads a Word document from a file and returns the text as a string.
func ReadFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Read(data)
}

// Read reads a turn report from a byte slice and returns the contents as a slice of strings.
func Read(data []byte) ([]string, error) {
	// normalize line endings before splitting.
	return strings.Split(string(norm.LineEndings(data)), "\n"), nil
}
