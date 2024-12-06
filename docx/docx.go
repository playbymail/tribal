// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package docx provides a reader for Microsoft Word documents.
//
// The reader is copied from https://github.com/lu4p/cat/blob/master/docxtxt/docxreader.go.
// All credit goes to the original author, https://github.com/lu4p.
//
// A few modifications were made to make it work with the Tribal reader.
// All errors are mine, naturally.
package docx

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
)

// http://officeopenxml.com/anatomyofOOXML.php

// ReadFile loads a Word document from a file and returns the text as a slice of byte slices.
// Each slice of bytes represents a paragraph in the original document.
func ReadFile(path string) ([][]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Read(data)
}

// Read reads a Word document from a byte slice and returns the contents as a slice of byte slices.
// Each slice of bytes is a single paragraph from the document.
func Read(data []byte) ([][]byte, error) {
	r := bytes.NewReader(data)

	// Word documents are zip files.
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, err
	}

	var doc docx

	// Find the document XML file in the ZIP archive.
	const docName = "word/document.xml"
	var docFile *zip.File
	for _, f := range zr.File {
		if f.Name == docName {
			docFile = f
			break
		}
	}
	if docFile == nil {
		return nil, errors.Join(ErrInvalidDocument, errors.New(docName+" file not found"))
	}

	// Extract the XML data from the ZIP file.
	var docXML []byte
	rdr, err := docFile.Open()
	defer rdr.Close()
	if err != nil {
		return nil, errors.Join(ErrInvalidDocument, err)
	} else if data, err := io.ReadAll(rdr); err != nil {
		return nil, errors.Join(ErrInvalidDocument, err)
	} else {
		docXML = data
	}

	// Tokenize the XML data into paragraphs and words.
	doc.listP(docXML)

	// Combine the tokens into a slice of slices containing all the paragraphs.
	paragraphs := make([][]byte, 0, len(doc.WordsList))
	for _, word := range doc.WordsList {
		paragraphs = append(paragraphs, bytes.Join(word.Words, []byte(" ")))
	}

	return paragraphs, nil
}

// docx struct for managing the Word document content.
type docx struct {
	WordsList []*words
}

type words struct {
	Words [][]byte
}

var (
	// rxParagraph matches an entire paragraph tag (`<w:p>...</w:p>`) and captures the content within it.
	// Breakdown:
	//  `(?U)`: Enables ungreedy mode, making `.*` match the shortest possible string.
	//  `<w:p[^>]*>`: Matches the opening `<w:p>` tag. It can include optional attributes (e.g., `<w:p attr="value">`).
	//  `<w:p`: Matches the start of the tag.
	//  `[^>]*`: Matches zero or more characters that are not `>`, to allow attributes inside the tag.
	//  `(.*)`: Captures everything between the opening `<w:p>` and closing `</w:p>` tags.
	//  `</w:p>`: Matches the closing tag.
	rxParagraph = regexp.MustCompile(`(?U)<w:p[^>]*>(.*)</w:p>`)

	// rxListP is similar to `rxParagraph`, but the placement of the capture group suggests a subtle difference in intent.
	// Difference:
	//   `(.*)` in `rxListP` starts capturing after `[^>]*` rather than after the opening `<w:p>` tag.
	// Practical Effect:
	//   Both expressions behave similarly in many cases, but `rxListP` may behave differently if additional
	//   attributes or nested content are present. This depends on the XML's structure.
	rxListP = regexp.MustCompile(`(?U)<w:p[^>]*(.*)</w:p>`)

	// rxRunOfText matches `<w:r>` tags, which are often used in Word documents to group runs of text, and captures the content between the tags.
	// Breakdown:
	//  `(?U)`: Enables ungreedy mode.
	//  `(<w:r>|<w:r .*>)`: Matches an opening `<w:r>` tag, with or without attributes.
	//  `<w:r>`: Matches a plain `<w:r>` tag.
	//  `<w:r .*>`: Matches a `<w:r>` tag with attributes (e.g., `<w:r attr="value">`).
	//  `|`: Allows either the plain tag or the tag with attributes and is captured as the first group.
	//  `(.*)`: Captures the content between the `<w:r>` and `</w:r>` tags and is captured as the second group.
	//  `(</w:r>)`: Matches the closing `</w:r>` tag and is captured as the third group.
	rxRunOfText = regexp.MustCompile(`(?U)(<w:r>|<w:r .*>)(.*)(</w:r>)`)

	// rxText matches `<w:t>` tags, which contain the actual text content in Word documents, and captures the text between the tags.
	// Breakdown:
	// ` (?U)`: Enables ungreedy mode.
	//  `(<w:t>|<w:t .*>)`: Matches an opening `<w:t>` tag, with or without attributes and is captured as the first group.
	//    `<w:t>`: Matches a plain `<w:t>` tag.
	//    `<w:t .*>`: Matches a `<w:t>` tag with attributes and is captured as the first group.
	//  `(.*)`: Captures the content between the `<w:t>` and `</w:t>` tags and is captured as the second group.
	//  `(</w:t>)`: Matches the closing `</w:t>` tag and is captured as the third group.
	rxText = regexp.MustCompile(`(?U)(<w:t>|<w:t .*>)(.*)(</w:t>)`)

	//---
	// Summary of Use Cases
	// **`rxP` and `rxListP`**: Used to identify and process paragraphs (`<w:p>` elements).
	// **`rxRunT`**: Used to process runs of text grouped under `<w:r>` tags.
	// **`rxT`**: Used to extract the actual text from `<w:t>` tags.
)

// listP parses paragraphs from the document XML.
// It recursively processes nested paragraphs and somehow updates the receiver.
func (d *docx) listP(data []byte) {
	for _, match := range rxListP.FindAllSubmatch(data, -1) {
		item := match[1]
		if rxParagraph.Match(item) {
			// the paragraph contains another paragraph
			d.listP(item)
		} else {
			// Extract text from w:t tags.
			var content [][]byte
			for _, rMatch := range rxRunOfText.FindAllSubmatch(item, -1) {
				for _, textMatch := range rxText.FindAllSubmatch(rMatch[2], -1) {
					content = append(content, textMatch[2])
				}
			}
			d.WordsList = append(d.WordsList, &words{Words: content})
		}
	}
}
