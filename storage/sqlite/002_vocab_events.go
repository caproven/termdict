package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/caproven/termdict/vocab"
	"github.com/oklog/ulid/v2"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upVocabEvents, downVocabEvents)
}

func upVocabEvents(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `CREATE TABLE vocab_events (
		id        TEXT    NOT NULL PRIMARY KEY,
		type      TEXT    NOT NULL,
		word      TEXT    NOT NULL COLLATE nocase,
		timestamp INTEGER NOT NULL
	)`); err != nil {
		return fmt.Errorf("create vocab events table: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `SELECT word, creation_timestamp FROM vocab`)
	if err != nil {
		return fmt.Errorf("query existing vocab: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Warn("Failed to close rows", "error", err)
		}
	}()

	for rows.Next() {
		var word string
		var ts int64
		if err := rows.Scan(&word, &ts); err != nil {
			return fmt.Errorf("scan vocab entry: %w", err)
		}
		id := ulid.MustNewDefault(time.Unix(ts, 0))
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO vocab_events (id, type, word, timestamp) VALUES (?, ?, ?, ?)`,
			id.String(), string(vocab.EventTypeAdd), word, ts); err != nil {
			return fmt.Errorf("insert vocab event: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("scan vocab entries: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DROP TABLE vocab`); err != nil {
		return fmt.Errorf("drop existing vocab table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `CREATE TABLE vocab (
		word TEXT NOT NULL PRIMARY KEY COLLATE nocase
	)`); err != nil {
		return fmt.Errorf("recreate vocab table: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO vocab (word) SELECT word FROM vocab_events WHERE type = 'add'`); err != nil {
		return fmt.Errorf("insert vocab entry: %w", err)
	}

	return nil
}

func downVocabEvents(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(ctx, `SELECT word FROM vocab`)
	if err != nil {
		return fmt.Errorf("query vocab: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Warn("Failed to close rows", "error", err)
		}
	}()

	var words []string
	for rows.Next() {
		var w string
		if err := rows.Scan(&w); err != nil {
			return fmt.Errorf("scan vocab entry: %w", err)
		}
		words = append(words, w)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("query vocab: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DROP TABLE vocab_events`); err != nil {
		return fmt.Errorf("drop vocab events table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DROP TABLE vocab`); err != nil {
		return fmt.Errorf("drop vocab table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `CREATE TABLE vocab (
		id                 INTEGER PRIMARY KEY,
		word               TEXT    NOT NULL UNIQUE COLLATE nocase,
		creation_timestamp INTEGER NOT NULL DEFAULT (unixepoch())
	)`); err != nil {
		return fmt.Errorf("recreate vocab table: %w", err)
	}

	for _, w := range words {
		if _, err := tx.ExecContext(ctx, `INSERT INTO vocab (word) VALUES (?)`, w); err != nil {
			return fmt.Errorf("insert vocab entry: %w", err)
		}
	}

	return nil
}
