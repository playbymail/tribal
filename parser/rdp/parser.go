// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package rdp implements a recursive descent parser for TribeNet turn reports.
package rdp

import (
	"os"
)

type Node any

func ParseAlloc(input []byte) *Parser {
	p := &Parser{
		input: normalize(input),
	}
	p.length = len(p.input)
	return p
}

type Parser struct {
	input  []byte
	pos    int
	length int
}

func (p *Parser) Buffer() []byte {
	return p.input
}

func (p *Parser) Parse() []*SectionId {
	var sections []*SectionId
	for s := p.AcceptSectionId(); s != nil; s = p.AcceptSectionId() {
		sections = append(sections, s)
	}
	return sections
}

func (p *Parser) WriteBuffer(path string) error {
	return os.WriteFile(path, p.input, 0644)
}

func (p *Parser) iseof() bool {
	return p.pos >= p.length
}

type Section struct{}

// ParseSection parses an entire section.
func (p *Parser) ParseSection() *Section {
	return nil
}
