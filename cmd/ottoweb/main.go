// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Ltime)
	started := time.Now()
	if err := run(); err != nil {
		log.Fatal(err)
	}
	log.Printf("ottoweb: completed in %v", time.Since(started))
}

func run() error {
	return nil
}
