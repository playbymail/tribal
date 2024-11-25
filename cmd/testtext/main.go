// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
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
		if err := run(name); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("text: read files in %v", time.Since(started))
}

func run(name string) error {
	started := time.Now()
	data, err := text.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	for n, line := range data {
		log.Printf("%d: \"%s\"", n, line)
	}
	log.Printf("%s: %6d lines in %v", name, len(data), time.Since(started))
	return nil
}
