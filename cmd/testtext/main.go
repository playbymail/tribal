// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/playbymail/tribal/section"
	"github.com/playbymail/tribal/text"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Ltime)
	started := time.Now()
	for _, name := range []string{
		"0899-12.0138.report.txt",
		"0900-05.0138.report.txt",
	} {
		if err := run(name, false); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("text: read files in %v", time.Since(started))
}

func run(name string, printLines bool) error {
	started := time.Now()
	// load the file
	data, err := text.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s: %4d lines in %v", name, len(data), time.Since(started))
	// optionally, print the lines
	if printLines {
		for n, line := range data {
			log.Printf("text: %4d: \"%s\"", n, line)
		}
	}
	// parse the report text into sections
	sections := section.Split(data)
	log.Printf("%s: %4d sections in %v", name, len(sections), time.Since(started))
	return nil
}
