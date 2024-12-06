// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"github.com/playbymail/tribal/store"
	"github.com/spf13/cobra"
	"log"
)

var (
	argsCreate struct {
		database string // path to the database file
	}

	cmdCreate = &cobra.Command{
		Use: "create",
	}

	argsCreateClan struct {
		clanNo int
	}

	cmdCreateClan = &cobra.Command{
		Use:   "clan",
		Short: "add a new clan to the  database",
		Long:  "add a new clan to the database",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if argsCreate.database == "" {
				return errors.New("database is required")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			s, err := store.Open(argsCreate.database, context.Background())
			if err != nil {
				log.Fatalf("error opening database: %v", err)
			}

			if err := runCreateClan(s, argsCreateClan.clanNo); err != nil {
				log.Fatalf("error creating clan: %v", err)
			}

			log.Printf("create: clan: %d: created", argsCreateClan.clanNo)
		},
	}
)

func runCreateClan(s *store.Store, clanNo int) error {
	if !(1 <= clanNo && clanNo <= 999) {
		return store.ErrInvalidClanId
	}
	if id, err := s.CreateClan(clanNo); err != nil {
		return err
	} else if id != clanNo {
		return errors.New("clan id is not the same as the one provided")
	}
	return nil
}
