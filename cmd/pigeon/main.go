// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the pigeon command.
package main

import (
	"github.com/playbymail/tribal/parser/pigeon"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	started := time.Now()
	defer func() {
		log.Printf("pigeon: finished in %v\n", time.Since(started))
	}()

	if err := run("9999-12.xxxx.report.txt"); err != nil {
		log.Fatal(err)
	}
}

func run(path string) error {
	started := time.Now()
	defer func() {
		log.Printf("pigeon: finished in %v\n", time.Since(started))
	}()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	log.Printf("pigeon: read %d bytes", len(data))

	return doit(data)
}

func doit(data []byte) error {
	started := time.Now()
	var px *pigeon.Parser
	defer func() {
		log.Printf("pigeon: doit completed in %v\n", time.Since(started))
	}()

	var err error
	px, err = pigeon.New(pigeon.WithData(data))
	if err != nil {
		return err
	}

	_, err = px.Parse()
	if err != nil {
		return err
	}

	return nil
}
