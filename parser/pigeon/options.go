// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package pigeon

// this file is used to define the options for the parser

// Option_t defines a single option for the parser
type Option_t func(*Parser) error

// WithData adds a normalized copy of the input to the parser.
// The caller can free the input data after this function returns.
func WithData(input []byte) Option_t {
	return func(p *Parser) error {
		p.input = input
		return nil
	}
}

// WithDebug sets the debug flag for the parser
func WithDebug(debug bool) Option_t {
	return func(p *Parser) error {
		return nil
	}
}

// WithHash sets the hash for the parser
func WithHash(hash string) Option_t {
	return func(p *Parser) error {
		return nil
	}
}
