// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/playbymail/tribal/parser/rdp"
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
	log.Printf("rdp: completed in %v", time.Since(started))
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

	now = time.Now()
	p := rdp.ParseAlloc(data)
	log.Printf("allocated parser           %8d bytes in %v", len(p.Buffer()), time.Since(now))

	now = time.Now()
	err = p.WriteBuffer("9999-12.xxxx.normalized.txt")
	if err != nil {
		return err
	}
	log.Printf("wrote                      normalized in %v", time.Since(now))

	now = time.Now()
	sections := p.Parse()
	for _, section := range sections {
		log.Printf("section: %+v\n", *section)
	}
	log.Printf("read                       %8d sections in %v", len(sections), time.Since(now))

	return nil
}
