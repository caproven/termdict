-- +goose Up
CREATE TABLE IF NOT EXISTS words
(
    id   INTEGER PRIMARY KEY,
    word TEXT NOT NULL UNIQUE COLLATE nocase
);

CREATE TABLE IF NOT EXISTS definitions
(
    id             INTEGER PRIMARY KEY,
    word_id        INTEGER NOT NULL,
    definition     TEXT    NOT NULL,
    part_of_speech TEXT    NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words (id) ON DELETE CASCADE
    -- TODO ON UPDATE? (restrict, cascade, ??)
);

CREATE INDEX IF NOT EXISTS idx_definitions_word_id ON definitions (word_id);
