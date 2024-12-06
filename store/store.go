// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package store implements the database access layer for the application.
// It provides methods for interacting with the database and performing CRUD operations.
// The only database supported is SQLite. We use the sqlc package to generate the code.
package store

//go:generate sqlc generate

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"github.com/playbymail/tribal"
	"github.com/playbymail/tribal/adapters"
	"github.com/playbymail/tribal/stdlib"
	"github.com/playbymail/tribal/store/sqlc"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	"time"
)

type Store struct {
	db  *sql.DB
	dbc *sqlc.Queries
	ctx context.Context
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// Open opens the database at the given path.
// If the database does not exist, it returns an error.
// The caller must call Close when done with the store.
func Open(path string, ctx context.Context) (*Store, error) {
	if ok, err := stdlib.IsFileExists(path); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrNotExist
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return &Store{db: db, dbc: sqlc.New(db), ctx: ctx}, nil
}

// CreateClan creates a new clan in the database.
func (s *Store) CreateClan(clanNo int) (int, error) {
	if !(1 <= clanNo && clanNo <= 999) {
		return 0, ErrInvalidClanId
	}
	err := s.dbc.CreateClan(s.ctx, sqlc.CreateClanParams{
		ID:   int64(clanNo),
		Name: fmt.Sprintf("%04d", clanNo),
	})
	if err != nil {
		log.Printf("insert failed: %v", err)
		if strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed: clans.id ") {
			return 0, ErrDuplicateClanId
		}
		return 0, errors.Join(ErrDatabase, err)
	}
	return clanNo, nil
}

// CreateReport loads the report for the given turn.
func (s *Store) CreateReport(rpt *tribal.ReportFile_t) error {
	// we never trust the client, so validate the input
	if rpt == nil {
		return errors.Join(ErrNoData, fmt.Errorf("report is nil"))
	} else if rpt.Hash == "" {
		return errors.Join(ErrNoData, fmt.Errorf("report is missing hash"))
	}
	year, month, ok := adapters.TurnIdToYearMonth(rpt.Turn)
	if !ok {
		return errors.Join(ErrInvalidTurnNo, fmt.Errorf("%d: invalid turn", rpt.Turn))
	}
	_, _ = year, month
	return ErrNotImplemented
}

// CreateTurn creates a new turn in the database.
// If the turn already exists, it ignores the request.
// Returns the turn ID or an error.
//
// Note: adding `ON CONFLICT (year, month) DO NOTHING` to the query
// will prevent the query from failing if the turn already exists.
func (s *Store) CreateTurn(year, month int) (int, error) {
	if !(899 <= year && year <= 9999) {
		return 0, ErrInvalidYear
	} else if !(1 <= month && month <= 12) {
		return 0, ErrInvalidMonth
	}
	turnNo := (year-899)*12 + month - 12
	err := s.dbc.CreateTurn(s.ctx, sqlc.CreateTurnParams{
		ID:    int64(turnNo),
		Year:  int64(year),
		Month: int64(month),
	})
	if err != nil {
		log.Printf("insert failed: %v", err)
		if !strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed: turns.id ") {
			return 0, errors.Join(ErrDatabase, err)
		}
	}
	return turnNo, nil
}

type ReportFileMeta_t struct {
	Id        int // key in the database
	Hash      string
	Name      string
	CreatedAt time.Time
}

// GetReportByHash returns the report for the given hash.
// Returns nil if the report does not exist.
func (s *Store) GetReportByHash(hash string) (*ReportFileMeta_t, error) {
	row, err := s.dbc.GetReportByHash(s.ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(ErrDatabase, err)
	}
	return &ReportFileMeta_t{
		Id:        int(row.ID),
		Hash:      hash,
		Name:      row.Name,
		CreatedAt: time.Unix(row.CreatedAt, 0).UTC(),
	}, nil
}

// GetTurnNo returns the turn number for the given year and month.
// If the turn does not exist, it returns an error.
func (s *Store) GetTurnNo(year, month int) (int, error) {
	if !(899 <= year && year <= 9999) {
		return 0, ErrInvalidYear
	} else if !(1 <= month && month <= 12) {
		return 0, ErrInvalidMonth
	}
	turnNo, err := s.dbc.GetTurnNo(s.ctx, sqlc.GetTurnNoParams{
		Year:  int64(year),
		Month: int64(month),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, errors.Join(ErrDatabase, err)
	}
	return int(turnNo), nil
}

// Hash returns the SHA1 hash of the given data.
func Hash(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}
