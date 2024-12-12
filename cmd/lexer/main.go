// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/playbymail/tribal/parser/lexer"
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
	log.Printf("read %d bytes", len(data))
	if len(data) != 0 && data[len(data)-1] != '\n' {
		data = append(data, '\n')
		log.Printf("read %d bytes", len(data))
	}

	lx, err := lexer.New(data)
	if err != nil {
		log.Fatal(err)
	}
	lexemes, lines := 0, 0
	for lexeme := lx.Next(); lexeme.Kind != lexer.EOF; lexeme = lx.Next() {
		lexemes++
		if lexeme.Kind == lexer.EOL {
			lines++
		}
		if lexemes < 55_555 {
			log.Printf("%7d:%7d: kind %2d: offset %8d: length %4d: %q\n", lines+1, lexemes, lexeme.Kind, lexeme.Offset, lexeme.Length, data[lexeme.Offset:lexeme.Offset+lexeme.Length])
		}
	}

	log.Printf("lexer: read %d lexemes, %d lines in %v", lexemes, lines, time.Since(started))
}
