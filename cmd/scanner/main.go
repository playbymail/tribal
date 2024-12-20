// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"github.com/playbymail/tribal/norm"
	"github.com/playbymail/tribal/parser/scanner"
	"github.com/playbymail/tribal/parser/units"
	"github.com/playbymail/tribal/section"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	started := time.Now()
	if err := run("9999-12.xxxx.report.txt"); err != nil {
		log.Fatal(err)
	}
	log.Printf("ottoweb: completed in %v", time.Since(started))
}

func orun() {
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
	log.Printf("scanner: read %d tokens in %v", s.NumTokens(), time.Since(started))

	started = time.Now()
	buf := &bytes.Buffer{}
	lines, priorTokenKind := 0, scanner.EOF
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

func run(path string) error {
	started := time.Now()
	defer log.Printf("completed in %v", time.Since(started))

	// load the consolidated report data
	now := time.Now()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	log.Printf("read                       %8d bytes in %v", len(data), time.Since(now))

	// normalize the line endings
	now = time.Now()
	data = norm.LineEndings(data)
	log.Printf("normalized line endings of %8d bytes in %v", len(data), time.Since(now))

	// split the data into lines
	now = time.Now()
	lines := bytes.Split(data, []byte{'\n'})
	log.Printf("split                      %8d bytes into %8d lines    in %v", len(data), len(lines), time.Since(now))

	// split the lines into sections
	now = time.Now()
	sections := section.Split(lines)
	log.Printf("sectioned                  %8d lines into %8d sections in %v", len(lines), len(sections), time.Since(now))

	now = time.Now()
	err = dumpSections(sections, "9999-12.xxxx.sections.txt")
	if err != nil {
		return err
	}
	log.Printf("dumped                     %8d sections in %v", len(sections), time.Since(now))

	// parse the sections, returning the map and all errors
	now = time.Now()
	for n, s := range sections {
		uht := units.ParseUnitHeading("path", s.Header)
		if uht == nil {
			log.Printf("not a unit heading: %q", s.Header)
		} else {
			log.Printf("unit heading: %+v", *uht)
		}
		if n > 3 {
			break
		}
	}
	log.Printf("parsed                     %8d sections in %v", len(sections), time.Since(now))

	return nil
}

func dumpSections(sections []*section.Section, path string) error {
	return os.WriteFile(path, section.DumpSections(sections), 0644)
}
