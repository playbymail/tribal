// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package rdp_test

import (
	"bytes"
	"github.com/playbymail/tribal/parser/rdp"
	"testing"
)

func TestAcceptToLF(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantTok  []byte
		wantRest []byte
	}{
		{
			name:     "empty buffer",
			input:    []byte{},
			wantTok:  nil,
			wantRest: []byte{},
		},
		{
			name:     "single LF",
			input:    []byte{rdp.LF},
			wantTok:  []byte{},
			wantRest: []byte{rdp.LF},
		},
		{
			name:     "text with LF",
			input:    []byte("hello\nworld"),
			wantTok:  []byte("hello"),
			wantRest: []byte{'\n', 'w', 'o', 'r', 'l', 'd'},
		},
		{
			name:     "no LF in text",
			input:    []byte("hello world"),
			wantTok:  nil,
			wantRest: []byte("hello world"),
		},
		{
			name:     "multiple LFs",
			input:    []byte("hello\nworld\n"),
			wantTok:  []byte("hello"),
			wantRest: []byte{'\n', 'w', 'o', 'r', 'l', 'd', '\n'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTok, gotRest := rdp.AcceptToLF(tt.input)
			if !bytes.Equal(gotTok, tt.wantTok) {
				t.Errorf("token = %q, want %q", gotTok, tt.wantTok)
			}
			if !bytes.Equal(gotRest, tt.wantRest) {
				t.Errorf("rest = %q, want %q", gotRest, tt.wantRest)
			}
		})
	}
}
