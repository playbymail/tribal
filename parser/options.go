// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package parser

// this file is used to define the options for the parser

type Options_t struct {
	Data  []byte // report input to parse
	Debug bool
	Hash  string // SHA1 hash of the report
}

type Option_t func(*Options_t) error

func WithData(data []byte) Option_t {
	return func(c *Options_t) error {
		c.Data = data
		return nil
	}
}

func WithDebug(debug bool) Option_t {
	return func(c *Options_t) error {
		c.Debug = debug
		return nil
	}
}

func WithHash(hash string) Option_t {
	return func(c *Options_t) error {
		c.Hash = hash
		return nil
	}
}
