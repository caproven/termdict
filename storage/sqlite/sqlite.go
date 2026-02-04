package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/caproven/termdict/dictionary"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Store struct {
	db *sql.DB
}

// NewStore constructs a store and performs db initialization.
func NewStore(ctx context.Context, db *sql.DB) (*Store, error) {
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		return nil, fmt.Errorf("enable foreign key constrains: %w", err)
	}

	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("set db migration dialect: %w", err)
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return nil, fmt.Errorf("apply db migrations: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) LookupWord(ctx context.Context, word string) ([]dictionary.Definition, error) {
	word = strings.ToLower(word)
	rows, err := s.db.QueryContext(ctx, `SELECT d.definition, d.part_of_speech FROM definitions AS d INNER JOIN words AS w ON d.word_id = w.id WHERE w.word IS ?`, word)
	if err != nil {
		return nil, fmt.Errorf("query definitions for word %q: %w", word, err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			slog.Warn("Failed to close rows", "error", err)
		}
	}(rows)

	var defs []dictionary.Definition
	for rows.Next() {
		var def dictionary.Definition
		if err := rows.Scan(&def.Meaning, &def.PartOfSpeech); err != nil {
			return nil, fmt.Errorf("scan definition for word %q: %w", word, err)
		}
		defs = append(defs, def)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query definitions for word %q: %w", word, rows.Err())
	}

	if len(defs) == 0 {
		return nil, fmt.Errorf("no definitions found for word %q", word)
	}

	return defs, nil
}

func (s *Store) ContainsWord(ctx context.Context, word string) (bool, error) {
	word = strings.ToLower(word)
	var exists int
	query := `SELECT EXISTS(SELECT 1 FROM words WHERE word = ?)`
	if err := s.db.QueryRowContext(ctx, query, word).Scan(&exists); err != nil {
		return false, fmt.Errorf("query word %q: %w", word, err)
	}

	return exists == 1, nil
}

func (s *Store) SaveWord(ctx context.Context, word string, defs []dictionary.Definition) (err error) {
	word = strings.ToLower(word)
	if len(strings.TrimSpace(word)) == 0 {
		return errors.New("word is blank")
	}
	if len(defs) == 0 {
		return errors.New("no definitions to save")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	res, err := tx.ExecContext(ctx, `INSERT INTO words (word) VALUES (?)`, word)
	if err != nil {
		return fmt.Errorf("insert word %q: %w", word, err)
	}
	wordID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last word id: %w", err)
	}

	defStatement, err := tx.PrepareContext(ctx, `INSERT INTO definitions (word_id, definition, part_of_speech) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare definition statement: %w", err)
	}
	for _, def := range defs {
		if _, err := defStatement.ExecContext(ctx, wordID, def.Meaning, def.PartOfSpeech); err != nil {
			return fmt.Errorf("insert definition for word %q: %w", word, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// AddWordsToList adds words to the list. Words that already exist are ignored, and newly inserted words are returned.
func (s *Store) AddWordsToList(ctx context.Context, words []string) ([]string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	var inserted []string

	insertStatement, err := tx.PrepareContext(ctx, `INSERT INTO vocab (word) VALUES (?) ON CONFLICT DO NOTHING`)
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %w", err)
	}
	for _, word := range words {
		word = strings.ToLower(word)
		res, err := insertStatement.ExecContext(ctx, word)
		if err != nil {
			return nil, fmt.Errorf("insert word %q: %w", word, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("get rows affected: %w", err)
		}
		if affected == 1 {
			inserted = append(inserted, word)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return inserted, nil
}

// RemoveWordsFromList removes words from the list. Words that don't exist are ignored, and words which are removed
// by this operation are returned.
func (s *Store) RemoveWordsFromList(ctx context.Context, words []string) ([]string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	var removed []string

	deleteStatement, err := tx.PrepareContext(ctx, `DELETE FROM vocab WHERE word = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %w", err)
	}
	for _, word := range words {
		word = strings.ToLower(word)
		res, err := deleteStatement.ExecContext(ctx, word)
		if err != nil {
			return nil, fmt.Errorf("remove word %q: %w", word, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("get rows affected: %w", err)
		}
		if affected == 1 {
			removed = append(removed, word)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return removed, nil
}

func (s *Store) GetWordsInList(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT word FROM vocab ORDER BY word`)
	if err != nil {
		return nil, fmt.Errorf("query words in list: %w", err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			slog.Warn("Failed to close rows", "error", err)
		}
	}(rows)

	var words []string
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, fmt.Errorf("scan word in list: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}
