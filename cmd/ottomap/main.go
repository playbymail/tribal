// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
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

	flag.BoolVar(&section.DebugConfig.SplitTurns, "split-turns", false, "Enable splitting of turns")
	flag.BoolVar(&section.DebugConfig.SplitFollows, "split-follows", false, "Enable splitting of follows")
	flag.BoolVar(&section.DebugConfig.SplitGoesTo, "split-goes-to", false, "Enable splitting of goesTo")
	flag.BoolVar(&section.DebugConfig.SplitMarches, "split-marches", false, "Enable splitting of marches")
	flag.BoolVar(&section.DebugConfig.SplitSails, "split-sails", false, "Enable splitting of sails")
	flag.BoolVar(&section.DebugConfig.SplitPatrols, "split-patrols", false, "Enable splitting of patrols")
	flag.BoolVar(&section.DebugConfig.SplitStatus, "split-status", false, "Enable splitting of status")

	flag.Parse()

	log.SetFlags(log.Lshortfile)

	var started time.Time

	log.Printf("ottomap: writing log file to %s\n", "ottomap.txt")
	if fd, err := os.OpenFile("ottomap.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(err)
	} else {
		log.SetOutput(fd)
	}
	started = time.Now()
	for _, clan := range []int{138} {
		for _, turnId := range []tribal.TurnId_t{0, 1, 2, 3, 4, 5} {
			if err := importWord(".", clan, turnId); err != nil {
				log.Print(err)
				log.SetOutput(os.Stderr)
				log.Fatal(err)
			}
		}
	}
	log.SetOutput(os.Stderr)

	log.Printf("ottomap: writing log file to %s\n", "0999-12.parser.txt")
	if fd, err := os.OpenFile("0999-12.parser.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(err)
	} else {
		log.SetOutput(fd)
	}
	started = time.Now()
	section.Units = nil
	if err := importPlainText(".", 999, 1200); err != nil {
		log.Print(err)
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
	log.Printf("ottomap: completed in %v", time.Since(started))
	log.SetOutput(os.Stderr)

	if buf, err := json.MarshalIndent(section.Units, "", "  "); err != nil {
		log.Fatal(err)
	} else if err = os.WriteFile("0999-12.parser.json", buf, 0644); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("ottomap: wrote 0999-12.parser.json in %v\n", time.Since(started))
	}
}

func importPlainText(path string, clan int, turnId tribal.TurnId_t) error {
	started := time.Now()
	defer log.Printf("import: plain text: completed in %v", time.Since(started))

	turnYear, turnMonth, ok := adapters.TurnIdToYearMonth(turnId)
	if !ok {
		return fmt.Errorf("invalid turn id: %d", turnId)
	}
	reportName := fmt.Sprintf("%04d-%02d.%04d.report.txt", turnYear, turnMonth, clan)
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
