// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package text

import (
	"bytes"
	"github.com/playbymail/tribal/norm"
	"os"
)

// ReadFile loads a Word document from a file and returns the text as a slice of lines.
func ReadFile(path string) ([][]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Read(data)
}

// Read reads a turn report from a byte slice and returns the contents as a slice of lines.
func Read(data []byte) ([][]byte, error) {
	// normalize line endings before splitting.
	return bytes.Split(norm.LineEndings(data), []byte{'\n'}), nil
}
