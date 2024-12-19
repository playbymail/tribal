// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package lemon

import (
	"fmt"
	"github.com/playbymail/tribal/parser/scanner"
	"strconv"
)

type Parser struct {
	s *scanner.Scanner
}

type Node struct {
	Type  Type
	Error error
}

func ParseAlloc(s *scanner.Scanner) *Parser {
	return &Parser{s: s}
}

func (p *Parser) Parse() (*Node, error) {
	return nil, nil
}

func (p *Parser) accept(types ...scanner.Type) (*scanner.Token, bool) {
	return p.s.Accept(types...)
}

func (p *Parser) expect(types ...scanner.Type) (*scanner.Token, error) {
	if token, ok := p.accept(types...); ok {
		return token, nil
	}
	token := p.s.Next()
	return token, fmt.Errorf("got %v", token.Type)
}

func (p *Parser) unit_section_list() (*Node, error) {
	return nil, nil
}

func (p *Parser) unit_section() (*Node, error) {
	return nil, nil
}

func (p *Parser) unit_header() (*Node, error) {
	return nil, nil
}

func (p *Parser) unit_header_id() (*Node, error) {
	if _, ok := p.accept(scanner.Courier); ok {
		token, err := p.expect(scanner.CourierID)
		if err != nil {
			// consume to end of line and return error
			p.s.RunTo(scanner.Newline, scanner.EOF)
			err = fmt.Errorf("expected courier id, got %q", token.String())
			return &Node{Type: ERROR}, err
		}
		return &Node{Type: UNIT_ID}, nil
	}
	if _, ok := p.accept(scanner.Element); ok {
		token, err := p.expect(scanner.ElementID)
		if err != nil {
			// consume to end of line and return error
			p.s.RunTo(scanner.Newline, scanner.EOF)
			err = fmt.Errorf("expected element id, got %q", token.String())
			return &Node{Type: ERROR}, err
		}
		return &Node{Type: UNIT_ID}, nil
	}
	if _, ok := p.accept(scanner.Fleet); ok {
		token, err := p.expect(scanner.FleetID)
		if err != nil {
			// consume to end of line and return error
			p.s.RunTo(scanner.Newline, scanner.EOF)
			err = fmt.Errorf("expected fleet id, got %q", token.String())
			return &Node{Type: ERROR}, err
		}
		return &Node{Type: UNIT_ID}, nil
	}
	if _, ok := p.accept(scanner.Garrison); ok {
		token, err := p.expect(scanner.GarrisonID)
		if err != nil {
			// consume to end of line and return error
			p.s.RunTo(scanner.Newline, scanner.EOF)
			return &Node{Type: ERROR}, err
		}
		return &Node{Type: UNIT_ID}, nil
	}
	if _, ok := p.accept(scanner.Tribe); ok {
		token, err := p.expect(scanner.Number)
		if err != nil {
			// consume to end of line and return error
			p.s.RunTo(scanner.Newline, scanner.EOF)
			err = fmt.Errorf("expected tribe id, got %q", token.String())
			return &Node{Type: ERROR}, err
		}
		n, err := strconv.Atoi(token.Value)
		if err != nil || !(0 < n && n <= 9999) {
			// consume to end of line and return error
			p.s.RunTo(scanner.Newline, scanner.EOF)
			err = fmt.Errorf("invalid tribe id %q", token.String())
			return &Node{Type: ERROR}, err
		}
		return &Node{Type: UNIT_ID}, nil
	}
	return &Node{Type: ERROR}, fmt.Errorf("expected unit id")
}

func (p *Parser) unit_name() (*Node, error) {
	_, ok := p.accept(scanner.Text)
	if ok {
		return &Node{Type: UNIT_NAME}, nil
	}
	return &Node{Type: EPSILON}, nil
}

type Type int

const (
	UNKNOWN Type = iota
	EPSILON
	ERROR
	UNIT_ID
	UNIT_NAME
)
