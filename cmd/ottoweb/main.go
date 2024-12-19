// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"github.com/playbymail/tribal/norm"
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
	log.Printf("split                      %8d lines into %8d sections in %v", len(lines), len(sections), time.Since(now))

	now = time.Now()
	err = dumpSections(sections, "9999-12.xxxx.sections.txt")
	if err != nil {
		return err
	}
	log.Printf("dumped                     %8d sections in %v", len(sections), time.Since(now))
	return nil
}

func dumpSections(sections []*section.Section, path string) error {
	b := &bytes.Buffer{}
	for _, s := range sections {
		b.Write(s.Header)
		b.WriteByte('\n')
		if len(s.Turn) != 0 {
			b.Write(s.Turn)
			b.WriteByte('\n')
		}
		if len(s.Moves.Movement) != 0 {
			b.Write(s.Moves.Movement)
			b.WriteByte('\n')
		}
		if len(s.Moves.Follows) != 0 {
			b.Write(s.Moves.Follows)
			b.WriteByte('\n')
		}
		if len(s.Moves.GoesTo) != 0 {
			b.Write(s.Moves.GoesTo)
			b.WriteByte('\n')
		}
		if len(s.Moves.Fleet) != 0 {
			b.Write(s.Moves.Fleet)
			b.WriteByte('\n')
		}
		if len(s.Moves.Scouts) != 0 {
			for _, scout := range s.Moves.Scouts {
				if len(scout) == 0 {
					continue
				}
				b.Write(scout)
				b.WriteByte('\n')
			}
		}
		if len(s.Status) != 0 {
			b.Write(s.Status)
			b.WriteByte('\n')
		}
	}
	return os.WriteFile(path, b.Bytes(), 0644)
}
