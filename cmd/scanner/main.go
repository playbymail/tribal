// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"github.com/playbymail/tribal/parser/scanner"
	"log"
	"os"
	"time"
)

func main() {
	started := time.Now()

	data, err := os.ReadFile("9999-12.xxxx.report.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("read %d bytes in %v", len(data), time.Since(started))

	started = time.Now()
	s := scanner.New(data)
	tokens, lines := len(s.Tokens()), 0
	log.Printf("scanner: read %d tokens, %d lines in %v", tokens, lines, time.Since(started))

	started = time.Now()
	buf := &bytes.Buffer{}
	priorTokenKind := scanner.EOF
	for token := s.Next(); token.Type != scanner.EOF; token = s.Next() {
		switch token.Type {
		case scanner.EOF:
			// do nothing
		case scanner.Newline:
			buf.WriteByte('\n')
			lines++
		case scanner.Text:
			buf.WriteString(fmt.Sprintf("%s", token.Value))
		case scanner.Unknown:
			// only print unknown tokens if we need to separate text
			if priorTokenKind == scanner.Text {
				buf.WriteByte('?')
			}
		case scanner.Whitespace:
			// avoid printing extra spaces
			if priorTokenKind != scanner.Whitespace {
				buf.WriteByte(' ')
			}
		}
		if token.Type == scanner.Unknown {
			priorTokenKind = scanner.Whitespace
		} else {
			priorTokenKind = token.Type

		}
	}
	log.Printf("scanner: buffered %d lines in %v", lines, time.Since(started))

	started = time.Now()
	if err := os.WriteFile("../output/9999-12.xxxx.scanned.txt", buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("scanner: wrote %d lines (%d bytes) in %v", lines, len(buf.Bytes()), time.Since(started))
}
