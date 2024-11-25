// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package docx provides a reader for Microsoft Word documents.
//
// The reader is copied from https://github.com/lu4p/cat/blob/master/docxtxt/docxreader.go.
// All credit goes to the original author, https://github.com/lu4p.
//
// A few modifications were made to make it work with the Tribal reader.
// All errors are mine, naturally.
package docx

// http://officeopenxml.com/anatomyofOOXML.php

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

// ReadFile loads a Word document from a file and returns the text as a slice of strings.
// Technically, each string is a paragraph. In practice, the turn report is composed of
// single line paragraphs.
func ReadFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Read(data)
}

// Read reads a Word document from a byte slice and returns the contents
// as a slice of strings. Each string is a single paragraph from the document.
func Read(data []byte) ([]string, error) {
	r := bytes.NewReader(data)

	// Word documents are zip files.
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, err
	}

	doc := docx{
		zipFileReader: nil,
		Files:         zr.File,
		FilesContent:  map[string][]byte{},
	}

	// the zip file will contain many files. we have to scan the
	// zip's directory listing to find the one we need.
	const docName = "word/document.xml"
	var docFolder *zip.File
	for _, f := range doc.Files {
		if f.Name == docName {
			docFolder = f
		}
	}
	if docFolder == nil {
		return nil, errors.New(docName + " file not found")
	}
	// extract the XML data from the zip folder
	var docXML string
	if rdr, err := docFolder.Open(); err != nil {
		return nil, err
	} else if data, err := io.ReadAll(rdr); err != nil {
		return nil, err
	} else {
		docXML = string(data)
	}

	// convert the xml data to tokens.
	// the tokens are stored as a list of paragraphs containing lists of words.
	// unfortunately, the tokenizer destroys Word tables.
	doc.listP(docXML)

	// combine the tokens into a slice containing all the paragraphs.
	// the docx library consume runs of spaces and tabs.
	// we compensate by injecting a space between every token in a paragraph.
	paragraphs := make([]string, 0, len(doc.WordsList))
	for _, word := range doc.WordsList {
		// word.Content holds the word tokens from a single paragraph.
		// we're going to save that as a single line in the result.
		paragraphs = append(paragraphs, strings.Join(word.Content, " "))
	}

	return paragraphs, nil
}

// docx zip struct
type docx struct {
	zipFileReader *zip.ReadCloser
	Files         []*zip.File
	FilesContent  map[string][]byte
	WordsList     []*words
}

type words struct {
	Content []string
}

// read all the data for a file in the zip file
func (d *docx) retrieveFileContents(filename string) ([]byte, error) {
	var file *zip.File
	for _, f := range d.Files {
		if f.Name == filename {
			file = f
		}
	}

	if file == nil {
		return nil, errors.New(filename + " file not found")
	}

	reader, err := file.Open()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(reader)
}

var (
	rxRunT = regexp.MustCompile(`(?U)(<w:r>|<w:r .*>)(.*)(</w:r>)`)
	rxT    = regexp.MustCompile(`(?U)(<w:t>|<w:t .*>)(.*)(</w:t>)`)
)

// get w:t value
func (d *docx) getT(item string) {
	var subStr string
	data := item
	w := new(words)
	content := []string{}

	wrMatch := rxRunT.FindAllStringSubmatchIndex(data, -1)
	// loop r
	for _, rMatch := range wrMatch {
		rData := data[rMatch[4]:rMatch[5]]
		wtMatch := rxT.FindAllStringSubmatchIndex(rData, -1)
		for _, match := range wtMatch {
			subStr = rData[match[4]:match[5]]
			content = append(content, subStr)
		}
	}
	w.Content = content
	d.WordsList = append(d.WordsList, w)
}

var (
	rxP = regexp.MustCompile(`(?U)<w:p[^>]*>(.*)</w:p>`)
)

// hasP identify the paragraph
func hasP(data string) bool {
	result := rxP.MatchString(data)
	return result
}

var (
	rxListP = regexp.MustCompile(`(?U)<w:p[^>]*(.*)</w:p>`)
)

// listP for w:p tag value
func (d *docx) listP(data string) {
	var result []string
	for _, match := range rxListP.FindAllStringSubmatch(data, -1) {
		result = append(result, match[1])
	}
	for _, item := range result {
		if hasP(item) {
			d.listP(item)
			continue
		}
		d.getT(item)
	}
}
