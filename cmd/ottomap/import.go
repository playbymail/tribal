// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/playbymail/tribal"
	"github.com/playbymail/tribal/adapters"
	"github.com/playbymail/tribal/parser"
	"github.com/playbymail/tribal/stdlib"
	"github.com/playbymail/tribal/store"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	argsImport struct {
		database string // path to the database file
	}

	cmdImport = &cobra.Command{
		Use: "import",
	}

	argsImportReport struct {
		clan int    // clan that owns the report
		path string // path to the report file
	}

	cmdImportReport = &cobra.Command{
		Use:   "report",
		Short: "import report into the database",
		Long:  "import report into the database",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if argsImport.database == "" {
				return errors.New("database is required")
			} else if !(0 <= argsImportReport.clan && argsImportReport.clan < 1000) {
				return errors.New("clan must be between 0 and 999")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			started := time.Now()
			if ok, err := stdlib.IsFileExists(argsImportReport.path); err != nil {
				log.Fatalf("error checking for file: %v", err)
			} else if !ok {
				log.Fatalf("file does not exist: %s", argsImportReport.path)
			}

			// check that the file is a report that matches the YYYY-MM.CLAN.report.(docx|txt) format.
			clan, turn, ok := adapters.ReportFileNameToClanTurn(argsImportReport.path)
			if !ok {
				log.Fatalf("import: invalid report name: %s", argsImportReport.path)
			}

			// use the clan from the command line if provided.
			// otherwise, use the clan from the file name.
			// this is a hack to support the GM who will import reports from multiple clans.
			if argsImportReport.clan != 0 {
				log.Printf("import: clan: overriding %d: %d\n", clan, argsImportReport.clan)
				clan, ok = adapters.IntToClanId(argsImportReport.clan)
				if !ok {
					log.Fatalf("import: invalid clan: %d", argsImportReport.clan)
				}
			}
			year, month := turn.YearMonth()
			log.Printf("import: clan %04d: turn %04d-%02d (#%d)\n", clan, year, month, turn)

			s, err := store.Open(argsImport.database, context.Background())
			if err != nil {
				log.Fatalf("error opening database: %v", err)
			}
			defer s.Close()

			if err := runImportReport(s, clan, turn, argsImportReport.path, true); err != nil {
				log.Fatalf("error importing report: %v", err)
			}
			log.Printf("import: report: %s: done in %v\n", argsImportReport.path, time.Since(started))
		},
	}
)

// runImportReport imports a report into the database.
// It returns an error if the report is not unique.
// It should update the database with a single transaction.
func runImportReport(s *store.Store, clan tribal.ClanId_t, turn tribal.TurnId_t, path string, debug bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	hash := store.Hash(data)
	if row, err := s.GetReportByHash(hash); err != nil {
		log.Printf("error: %v\n", err)
	} else if row != nil {
		log.Printf("error: report with same hash already exists\n")
		log.Printf("       name: %s\n", row.Name)
		log.Printf("       created at: %s\n", row.CreatedAt.Format(time.RFC3339))
		return store.ErrDuplicateReport
	}
	log.Printf("import: report: %s: seems unique\n", path)

	rpt, err := parser.Report(path, parser.WithDebug(debug), parser.WithData(data))
	if err != nil {
		return errors.Join(fmt.Errorf("parsing error"), err)
	} else if rpt == nil {
		log.Printf("import: report: %s: parser returned nil\n", path)
		return nil
	}
	log.Printf("import: report: %s: %s\n", path, rpt.Hash)
	if rpt.Turn == nil {
		rpt.Turn = &parser.Turn_t{
			No:    int(turn),
			Error: fmt.Errorf("report has no turn"),
		}
		rpt.Turn.Year, rpt.Turn.Month, _ = adapters.TurnIdToYearMonth(turn)
	} else if rpt.Turn.No != int(turn) {
		rpt.Turn.Error = fmt.Errorf("report has turn %d, expected %d", rpt.Turn.No, turn)
	}
	log.Printf("import: report: %s: %d units\n", path, len(rpt.Units))
	for _, unit := range rpt.Units {
		log.Printf("import: report: %s: unit %s\n", path, unit.Id)
	}

	// adapt from parser report to domain report
	drpt := tribal.ReportFile_t{
		Owner: clan,
		Turn:  turn,
		Hash:  rpt.Hash,
	}

	// this should be committed as a single transaction
	if err := s.CreateReport(&drpt); err != nil {
		return err
	}

	return nil
}
