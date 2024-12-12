// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
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
	// showDistributionOfBytes(data)

	started = time.Now()
	s, err := scanner.New(data, "899-12")
	if err != nil {
		log.Fatal(err)
	}
	tokens, lines := len(s.Tokens()), 0
	log.Printf("scanner: read %d tokens, %d lines in %v", tokens, lines, time.Since(started))

	started = time.Now()
	buf := &bytes.Buffer{}
	priorTokenKind := scanner.EOF
	for token := s.Next(); token.Type != scanner.EOF; token = s.Next() {
		if token.Type == scanner.Newline {
			lines++
			buf.WriteString(token.String())
			priorTokenKind = token.Type
		} else if token.Type == scanner.Unknown {
			// only print unknown tokens if we need to separate text
			if priorTokenKind == scanner.Text {
				buf.WriteString(token.String())
			}
			priorTokenKind = scanner.Whitespace
		} else if token.Type == scanner.Whitespace {
			// avoid printing extra spaces
			if priorTokenKind != scanner.Whitespace {
				buf.WriteString(token.String())
			}
			priorTokenKind = token.Type
		} else {
			buf.WriteString(token.String())
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

func showDistributionOfBytes(input []byte) {
	// show the distribution of bytes
	var tbl [256]int
	for _, ch := range input {
		tbl[ch]++
	}
	for i := 0; i < 256; i++ {
		if tbl[i] > 0 {
			log.Printf("%02x %c %d", i, i, tbl[i])
		}
	}
	if tbl['m'] != 0 {
		log.Fatalf("m not zero")
	}
}
