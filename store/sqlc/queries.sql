--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-- --------------------------------------------------------------------------
-- CreateClan creates a new clan.
--
-- name: CreateClan :exec
INSERT INTO clans (id, name)
VALUES (:id, :name);

-- --------------------------------------------------------------------------
-- CreateTurn creates a new turn.
-- If the turn already exists, it ignores the request.
--
-- name: CreateTurn :exec
INSERT INTO turns (id, year, month)
VALUES (:id, :year, :month);

-- --------------------------------------------------------------------------
-- GetReportByHash returns the report file with the given hash.
--
-- name: GetReportByHash :one
SELECT id, name, created_at
FROM report_files
WHERE hash = :hash;

-- --------------------------------------------------------------------------
-- GetTurnNo returns the turn number for the given year and month.
--
-- name: GetTurnNo :one
SELECT id
FROM turns
WHERE year = :year
  AND month = :month;