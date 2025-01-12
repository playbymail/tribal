// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the lemon command.
package main

import (
	"github.com/playbymail/tribal/parser/lemon"
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
	log.Printf("lemon: completed in %v", time.Since(started))
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

	// create a new parser with empty state and a tokenizer using our data
	now = time.Now()
	tokenizer := lemon.NewTokenizer(data)
	log.Printf("normalized input                      in %v\n", time.Since(now))

	now = time.Now()
	if err := tokenizer.WriteTo("9999-12.xxxx.normalized.txt"); err != nil {
		return err
	} else {
		log.Printf("wrote                      normalized in %v", time.Since(now))
	}

	parser := lemon.NewParser(path, data, tokenizer)

	parser.Parse()

	return nil
}
