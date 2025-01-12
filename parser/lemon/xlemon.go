// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package lemon

import (
	"fmt"
	"github.com/playbymail/tribal/parser/scanner"
	"github.com/playbymail/tribal/section"
)

func ParseAlloc(s *section.Section) *XLemon_t {
	return &XLemon_t{s: s}
}

func (p *XLemon_t) Parse() (*Node, error) {
	return &Node{
		Type:  ERROR,
		Error: ErrNotImplemented,
	}, nil
}

func (p *XLemon_t) accept(types ...scanner.Type) (*scanner.Token, bool) {
	panic("!")
}

func (p *XLemon_t) expect(types ...scanner.Type) (*scanner.Token, error) {
	panic("!")
}

func (p *XLemon_t) unit_section_list() (*Node, error) {
	panic("!")
}

func (p *XLemon_t) unit_section() (*Node, error) {
	panic("!")
}

func (p *XLemon_t) unit_header() (*Node, error) {
	panic("!")
}

func (p *XLemon_t) unit_header_id() (*Node, error) {
	panic("!")
}

func (p *XLemon_t) unit_name() (*Node, error) {
	_, ok := p.accept(scanner.Text)
	if ok {
		return &Node{Type: UNIT_NAME}, nil
	}
	return &Node{Type: EPSILON}, nil
}

type XLemon_t struct {
	s *section.Section
}

type Node struct {
	Type  Type
	Error error
}

type Type int

const (
	UNKNOWN Type = iota
	EPSILON
	ERROR
	UNIT_ID
	UNIT_NAME
)

func (t Type) String() string {
	switch t {
	case UNKNOWN:
		return "unknown"
	case EPSILON:
		return "epsilon"
	case ERROR:
		return "error"
	case UNIT_ID:
		return "unit_id"
	case UNIT_NAME:
		return "unit_name"
	}
	panic(fmt.Sprintf("assert(type != %d)", t))
}
