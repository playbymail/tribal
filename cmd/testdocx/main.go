// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/playbymail/tribal/docx"
	"github.com/playbymail/tribal/section"
	"github.com/playbymail/tribal/text"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Ltime)
	started := time.Now()
	if err := runDocxLoad(false); err != nil {
		log.Fatal(err)
	}
	if err := runTextLoad(false); err != nil {
		log.Fatal(err)
	}
	log.Printf("ottomap: completed in %v", time.Since(started))
}

func runDocxLoad(printLines bool) error {
	started := time.Now()
	for _, name := range []string{
		"0899-12.0138.report.docx",
		"0900-05.0138.report.docx",
	} {
		if err := loadDocxFiles(name, false); err != nil {
			return err
		}
	}
	log.Printf("docx: read files in %v", time.Since(started))
	return nil
}

func loadDocxFiles(name string, printLines bool) error {
	started := time.Now()
	// load the file
	data, err := docx.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s: %4d lines in %v", name, len(data), time.Since(started))
	// optionally, print the lines
	if printLines {
		for n, line := range data {
			log.Printf("docx: %4d: \"%s\"", n, line)
		}
	}
	// parse the report text into sections
	sections := section.Split(data)
	log.Printf("%s: %4d sections in %v", name, len(sections), time.Since(started))
	for n, ss := range sections {
		log.Printf("docx: %4d: %s", n, ss.Header)
	}
	return nil
}

func runTextLoad(printLines bool) error {
	started := time.Now()
	for _, name := range []string{
		"0899-12.0138.report.txt",
		"0900-05.0138.report.txt",
	} {
		if err := loadTextFiles(name, false); err != nil {
			return err
		}
	}
	log.Printf("text: read files in %v", time.Since(started))
	return nil
}

func loadTextFiles(name string, printLines bool) error {
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
			log.Printf("docx: %4d: \"%s\"", n, line)
		}
	}
	// parse the report text into sections
	sections := section.Split(data)
	log.Printf("%s: %4d sections in %v", name, len(sections), time.Since(started))
	for n, ss := range sections {
		log.Printf("text: %4d: %s", n, ss.Header)
	}
	return nil
}
