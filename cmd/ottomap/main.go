// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/spf13/cobra"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

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
}

var (
	cmdRoot = &cobra.Command{
		Use:   "ottomap",
		Short: "ottomap is a tool for managing tribal data",
		Long:  "ottomap is a tool for managing tribal data",
	}
)
