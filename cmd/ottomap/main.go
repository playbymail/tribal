// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"github.com/playbymail/tribal"
	"github.com/playbymail/tribal/adapters"
	"github.com/playbymail/tribal/docx"
	"github.com/playbymail/tribal/section"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	const logFileName = "ottomap.txt"
	log.Printf("ottomap: writing log file to %s\n", logFileName)

	if fd, err := os.OpenFile(logFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(err)
	} else {
		log.SetOutput(fd)
	}

	started := time.Now()
	for _, clan := range []int{138} {
		for _, turnId := range []tribal.TurnId_t{0, 1, 2, 3, 4, 5} {
			if err := importWord(".", clan, turnId); err != nil {
				log.Fatal(err)
			}
		}
	}

	log.SetOutput(os.Stderr)
	log.Printf("ottomap: completed in %v", time.Since(started))
}

func importPlainText(path string, clan int, turnId tribal.TurnId_t) error {
	started := time.Now()
	defer log.Printf("import: plain text: completed in %v", time.Since(started))

	turnYear, turnMonth, ok := adapters.TurnIdToYearMonth(turnId)
	if !ok {
		return fmt.Errorf("invalid turn id: %d", turnId)
	}
	reportName := fmt.Sprintf("%04d-%02d.%04d.report.text", turnYear, turnMonth, clan)
	log.Printf("import: plain text: %s", reportName)

	// load the file
	input, err := os.ReadFile(filepath.Join(path, reportName))
	if err != nil {
		return err
	}

	return importReport(clan, turnId, input)
}

func importWord(path string, clan int, turnId tribal.TurnId_t) error {
	started := time.Now()
	defer log.Printf("import: word: completed in %v", time.Since(started))

	turnYear, turnMonth, ok := adapters.TurnIdToYearMonth(turnId)
	if !ok {
		return fmt.Errorf("invalid turn id: %d", turnId)
	}
	reportName := fmt.Sprintf("%04d-%02d.%04d.report.docx", turnYear, turnMonth, clan)
	log.Printf("import: word: %s", reportName)

	// load the file
	lines, err := docx.ReadFile(filepath.Join(path, reportName))
	if err != nil {
		return err
	}
	log.Printf("%s: %4d lines in %v", reportName, len(lines), time.Since(started))

	return importReport(clan, turnId, bytes.Join(lines, []byte{'\n'}))
}

func importReport(clan int, turnId tribal.TurnId_t, input []byte) error {
	started := time.Now()
	defer log.Printf("import: report: completed in %v", time.Since(started))

	turnYear, turnMonth, ok := adapters.TurnIdToYearMonth(turnId)
	if !ok {
		return fmt.Errorf("invalid turn id: %d", turnId)
	}
	reportId := fmt.Sprintf("%04d-%02d.%04d", turnYear, turnMonth, clan)

	log.Printf("import: %04d: %04d-%02d (#%d) bytes %d\n", clan, turnYear, turnMonth, turnId, len(input))

	// split the input into sections
	now := time.Now()
	sections := section.Split(input)
	log.Printf("sectioned                  %8d bytes into %8d sections in %v", len(input), len(sections), time.Since(now))

	now = time.Now()
	if err := dumpSections(sections, fmt.Sprintf("%04d-%02d.%04d.sections.txt", turnYear, turnMonth, clan), true); err != nil {
		return err
	}
	log.Printf("dumped                     %8d sections in %v", len(sections), time.Since(now))

	// parse the sections, returning the map and all errors
	now = time.Now()
	for _, s := range sections {
		if err := s.Parse(reportId); err != nil {
			log.Printf("section: %d: %s: %v", s.Id, s.Lines.Unit, err)
		}
	}
	log.Printf("parsed                     %8d sections in %v\n\n\n", len(sections), time.Since(now))

	return nil
}

func runCobra() error {
	cmdRoot.AddCommand(cmdCreate)
	cmdCreate.PersistentFlags().StringVarP(&argsCreate.database, "database", "D", "tribal.sqlite", "path to the database file")

	cmdCreate.AddCommand(cmdCreateClan)
	cmdCreateClan.Flags().IntVarP(&argsCreateClan.clanNo, "clan-id", "c", 0, "clan id to create")
	if err := cmdCreateClan.MarkFlagRequired("clan-id"); err != nil {
		log.Fatalf("create: clan: id: %v\n", err)
	}

	cmdRoot.AddCommand(cmdImport)
	cmdImport.PersistentFlags().StringVarP(&argsImport.database, "database", "D", "tribal.sqlite", "path to the database file")

	cmdImport.AddCommand(cmdImportReport)
	cmdImportReport.Flags().StringVarP(&argsImportReport.path, "file", "p", "", "path to the report file")
	if err := cmdImportReport.MarkFlagRequired("file"); err != nil {
		log.Fatalf("import: report: file: %v\n", err)
	}

	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}

	return nil
}

var (
	cmdRoot = &cobra.Command{
		Use:   "ottomap",
		Short: "ottomap is a tool for managing tribal data",
		Long:  "ottomap is a tool for managing tribal data",
	}
)

func dumpSections(sections []*section.Section, path string, separateUnits bool) error {
	return os.WriteFile(path, section.DumpSections(sections, separateUnits), 0644)
}
