package sqlite

import (
	"database/sql"
	"io"
	"testing"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) (_ *sql.DB) {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	return db
}

func closeAndAssertError(t testing.TB, closer io.Closer) {
	assert.NoError(t, closer.Close())
}

func TestNewStore(t *testing.T) {
	db := newTestDB(t)
	defer closeAndAssertError(t, db)
	_, err := NewStore(t.Context(), db)
	require.NoError(t, err)

	// Assert foreign key constraints enabled for connection
	var foreignKeys int
	require.NoError(t, db.QueryRowContext(t.Context(), `PRAGMA foreign_keys`).Scan(&foreignKeys))
	assert.Equal(t, 1, foreignKeys)

	// Check db setup occurred
	var wordsCount int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM words`).Scan(&wordsCount))
	assert.Equal(t, 0, wordsCount)

	var definitionsCount int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM definitions`).Scan(&definitionsCount))
	assert.Equal(t, 0, definitionsCount)

	// Check migrations and setup are idempotent
	_, err = NewStore(t.Context(), db)
	assert.NoError(t, err)

	// Check vocab_events table exists
	var eventsCount int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM vocab_events`).Scan(&eventsCount))
	assert.Equal(t, 0, eventsCount)

	// Check vocab table exists with new schema
	var vocabCount int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM vocab`).Scan(&vocabCount))
	assert.Equal(t, 0, vocabCount)
}

func TestStore_LookupWord(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO words (word) VALUES ('foo');
INSERT INTO definitions (word_id, definition, part_of_speech) VALUES
(last_insert_rowid(), 'def 1', 'noun'),
(last_insert_rowid(), 'def 2', 'adjective');`)
		require.NoError(t, err)

		defs, err := store.LookupWord(t.Context(), "foo")
		require.NoError(t, err)
		expectedDefs := []dictionary.Definition{
			{PartOfSpeech: "noun", Meaning: "def 1"},
			{PartOfSpeech: "adjective", Meaning: "def 2"},
		}
		assert.Equal(t, defs, expectedDefs)

		// check case ignored for lookup
		defs, err = store.LookupWord(t.Context(), "FOO")
		require.NoError(t, err)
		assert.Equal(t, expectedDefs, defs)
	})

	t.Run("word exists with no definitions", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO words (word) VALUES ('foo')`)
		require.NoError(t, err)

		defs, err := store.LookupWord(t.Context(), "foo")
		assert.Error(t, err)
		assert.Len(t, defs, 0)
	})
}

func TestStore_ContainsWord(t *testing.T) {
	db := newTestDB(t)
	defer closeAndAssertError(t, db)
	store, err := NewStore(t.Context(), db)
	require.NoError(t, err)

	_, err = db.ExecContext(t.Context(), `INSERT INTO words (word) VALUES ('foo')`)
	require.NoError(t, err)

	tests := map[string]struct {
		word    string
		want    bool
		wantErr bool
	}{
		"word exists matching case": {
			word: "foo",
			want: true,
		},
		"word exists non-matching case": {
			word: "FOO",
			want: true,
		},
		"word doesn't exist": {
			word: "bar",
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ok, err := store.ContainsWord(t.Context(), tt.word)
			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, ok)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, ok)
		})
	}
}

func TestStore_SaveWord(t *testing.T) {
	t.Run("lowercase word", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		defs := []dictionary.Definition{
			{PartOfSpeech: "verb", Meaning: "def 1"},
			{PartOfSpeech: "adverb", Meaning: "def 2"},
		}
		err = store.SaveWord(t.Context(), "foo", defs)
		require.NoError(t, err)

		gotDefs, err := store.LookupWord(t.Context(), "foo")
		require.NoError(t, err)
		assert.Equal(t, defs, gotDefs)

		// Ensure only a single row written to words table
		var wordCount int
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM words`).Scan(&wordCount))
		assert.Equal(t, 1, wordCount)

		// Ensure only the expected rows written to definitions table
		var defCount int
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM definitions`).Scan(&defCount))
		assert.Equal(t, 2, defCount)
	})

	t.Run("uppercase word transformed to lowercase", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		defs := []dictionary.Definition{
			{PartOfSpeech: "conjunction", Meaning: "def 1"},
		}
		err = store.SaveWord(t.Context(), "BAR", defs)
		require.NoError(t, err)

		gotDefs, err := store.LookupWord(t.Context(), "bar")
		require.NoError(t, err)
		assert.Equal(t, defs, gotDefs)
	})

	t.Run("invalid inputs", func(t *testing.T) {
		tests := map[string]struct {
			word string
			defs []dictionary.Definition
		}{
			"empty word": {
				word: "",
				defs: []dictionary.Definition{
					{PartOfSpeech: "verb", Meaning: "def 1"},
				},
			},
			"whitespace word": {
				word: " ",
				defs: []dictionary.Definition{
					{PartOfSpeech: "verb", Meaning: "def 1"},
				},
			},
			"no definitions": {
				word: "foo",
				defs: []dictionary.Definition{},
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				db := newTestDB(t)
				defer closeAndAssertError(t, db)
				store, err := NewStore(t.Context(), db)
				require.NoError(t, err)

				err = store.SaveWord(t.Context(), tt.word, tt.defs)
				assert.Error(t, err)

				// Ensure no entries written to words table
				var wordCount int
				require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM words`).Scan(&wordCount))
				assert.Equal(t, 0, wordCount)

				// Ensure no entries written to definitions table
				var defCount int
				require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count() FROM definitions`).Scan(&defCount))
				assert.Equal(t, 0, defCount)
			})
		}
	})

	t.Run("word already exists", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO words (word) VALUES ('foo')`)
		require.NoError(t, err)

		err = store.SaveWord(t.Context(), "foo", []dictionary.Definition{
			{PartOfSpeech: "verb", Meaning: "def 1"},
		})
		assert.Error(t, err)
		// TODO check ids or something
	})
}

func TestStore_AddWordsToList(t *testing.T) {
	t.Run("all new words", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		list := []string{"cascade", "dour"}
		added, err := store.AddWordsToList(t.Context(), list)
		require.NoError(t, err)
		assert.Equal(t, list, added)

		got := getVocabList(t, db)
		assert.Equal(t, list, got)
	})

	t.Run("new words with some existing", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('foo')`)
		require.NoError(t, err)

		added, err := store.AddWordsToList(t.Context(), []string{"foo", "bar", "baz"})
		require.NoError(t, err)
		assert.Equal(t, []string{"bar", "baz"}, added)

		got := getVocabList(t, db)
		assert.Equal(t, []string{"bar", "baz", "foo"}, got)
	})

	t.Run("all existing words", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('foo'), ('bar')`)
		require.NoError(t, err)

		added, err := store.AddWordsToList(t.Context(), []string{"foo", "bar"})
		require.NoError(t, err)
		assert.Len(t, added, 0)

		got := getVocabList(t, db)
		assert.Equal(t, []string{"bar", "foo"}, got)
	})

	t.Run("capitalization ignored for inserts", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		added, err := store.AddWordsToList(t.Context(), []string{"IRRESOLUTE"})
		require.NoError(t, err)
		assert.Equal(t, []string{"irresolute"}, added)

		got := getVocabList(t, db)
		assert.Equal(t, []string{"irresolute"}, got)
	})

	t.Run("capitalization ignored for duplicates", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('cacophony')`)
		require.NoError(t, err)

		added, err := store.AddWordsToList(t.Context(), []string{"CACOPHONY"})
		require.NoError(t, err)
		assert.Len(t, added, 0)

		got := getVocabList(t, db)
		assert.Equal(t, []string{"cacophony"}, got)
	})
}

func TestStore_RemoveWordsFromList(t *testing.T) {
	t.Run("all words to delete exist", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('tepid'), ('surmise'), ('eschew')`)
		require.NoError(t, err)

		// Out of order from inserts to show order doesn't matter
		removed, err := store.RemoveWordsFromList(t.Context(), []string{"eschew", "tepid", "surmise"})
		require.NoError(t, err)
		assert.Equal(t, []string{"eschew", "tepid", "surmise"}, removed)

		got := getVocabList(t, db)
		assert.Empty(t, got)
	})

	t.Run("delete subset of existing words", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('pyrrhic'), ('pervade')`)
		require.NoError(t, err)

		removed, err := store.RemoveWordsFromList(t.Context(), []string{"pyrrhic"})
		require.NoError(t, err)
		assert.Equal(t, []string{"pyrrhic"}, removed)

		got := getVocabList(t, db)
		assert.Equal(t, []string{"pervade"}, got)
	})

	t.Run("some words to delete exist", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('foo'), ('bar')`)
		require.NoError(t, err)

		removed, err := store.RemoveWordsFromList(t.Context(), []string{"foo", "baz"})
		require.NoError(t, err)
		assert.Equal(t, []string{"foo"}, removed)

		got := getVocabList(t, db)
		assert.Equal(t, []string{"bar"}, got)
	})

	t.Run("no words to delete exist", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('ambivalence')`)
		require.NoError(t, err)

		removed, err := store.RemoveWordsFromList(t.Context(), []string{"qwerty", "dvorak"})
		require.NoError(t, err)
		assert.Len(t, removed, 0)

		got := getVocabList(t, db)
		assert.Equal(t, []string{"ambivalence"}, got)
	})

	t.Run("capitalization ignored", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('herbaceous')`)
		require.NoError(t, err)

		removed, err := store.RemoveWordsFromList(t.Context(), []string{"HERBACEOUS"})
		require.NoError(t, err)
		assert.Equal(t, []string{"herbaceous"}, removed)

		got := getVocabList(t, db)
		assert.Empty(t, got)
	})
}

func TestStore_GetWordsInList(t *testing.T) {
	t.Run("no words", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		got, err := store.GetWordsInList(t.Context())
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("words alphabetically sorted at insertion", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('aardvark'), ('zebra')`)
		require.NoError(t, err)

		got, err := store.GetWordsInList(t.Context())
		require.NoError(t, err)
		assert.Equal(t, []string{"aardvark", "zebra"}, got)
	})

	t.Run("words not alphabetically sorted at insertion", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = db.ExecContext(t.Context(), `INSERT INTO vocab (word) VALUES ('zebra'), ('aardvark')`)
		require.NoError(t, err)

		got, err := store.GetWordsInList(t.Context())
		require.NoError(t, err)
		assert.Equal(t, []string{"aardvark", "zebra"}, got)
	})
}

func TestStore_GetEvents(t *testing.T) {
	t.Run("no events", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		events, err := store.GetEvents(t.Context())
		require.NoError(t, err)
		assert.Empty(t, events)
	})

	t.Run("returns events ordered by timestamp", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		_, err = store.AddWordsToList(t.Context(), []string{"bar", "foo"})
		require.NoError(t, err)

		events, err := store.GetEvents(t.Context())
		require.NoError(t, err)
		require.Len(t, events, 2)
		assert.Equal(t, vocab.EventTypeAdd, events[0].Type)
		assert.Equal(t, "bar", events[0].Word)
		assert.Equal(t, vocab.EventTypeAdd, events[1].Type)
		assert.Equal(t, "foo", events[1].Word)
		assert.NotEmpty(t, events[0].ID)
		assert.NotEmpty(t, events[1].ID)
		assert.True(t, events[0].Timestamp <= events[1].Timestamp, "events should be ordered by timestamp")
	})
}

func TestStore_AddEvents(t *testing.T) {
	t.Run("inserts new events and rebuilds vocab", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		events := []vocab.Event{
			{ID: "01J3XYZ1", Type: vocab.EventTypeAdd, Word: "foo", Timestamp: 100},
			{ID: "01J3XYZ2", Type: vocab.EventTypeAdd, Word: "bar", Timestamp: 200},
			{ID: "01J3XYZ3", Type: vocab.EventTypeRemove, Word: "foo", Timestamp: 300},
		}

		err = store.AddEvents(t.Context(), events)
		require.NoError(t, err)

		got, err := store.GetWordsInList(t.Context())
		require.NoError(t, err)
		assert.Equal(t, []string{"bar"}, got)

		stored, err := store.GetEvents(t.Context())
		require.NoError(t, err)
		assert.Len(t, stored, 3)
	})

	t.Run("ignores duplicate event IDs", func(t *testing.T) {
		db := newTestDB(t)
		defer closeAndAssertError(t, db)
		store, err := NewStore(t.Context(), db)
		require.NoError(t, err)

		events := []vocab.Event{
			{ID: "01J3XYZ1", Type: vocab.EventTypeAdd, Word: "foo", Timestamp: 100},
		}

		err = store.AddEvents(t.Context(), events)
		require.NoError(t, err)

		// Same ID, should be ignored
		events2 := []vocab.Event{
			{ID: "01J3XYZ1", Type: vocab.EventTypeAdd, Word: "foo", Timestamp: 100},
			{ID: "01J3XYZ2", Type: vocab.EventTypeAdd, Word: "bar", Timestamp: 200},
		}
		err = store.AddEvents(t.Context(), events2)
		require.NoError(t, err)

		stored, err := store.GetEvents(t.Context())
		require.NoError(t, err)
		assert.Len(t, stored, 2)
	})
}

func getVocabList(t testing.TB, db *sql.DB) []string {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), `SELECT word FROM vocab ORDER BY word`)
	require.NoError(t, err)
	var list []string
	for rows.Next() {
		var word string
		require.NoError(t, rows.Scan(&word))
		list = append(list, word)
	}
	return list
}
