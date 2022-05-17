package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/caproven/termdict/vocab"
)

const vocabFile string = "vocab.json"

// VocabRepo lets you fetch or save a vocab list
type VocabRepo struct {
	Path string
}

// Load a vocab list from storage
func (r VocabRepo) Load() (vocab.List, error) {
	var vl vocab.List

	if _, err := os.Stat(r.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// file doesn't exist, just return empty vocab.List
			return vl, nil
		}
		return vocab.List{}, err
	}

	data, err := os.ReadFile(r.Path)
	if err != nil {
		return vocab.List{}, err
	}

	if err := json.Unmarshal(data, &vl.Words); err != nil {
		return vocab.List{}, err
	}

	return vl, nil
}

// Save a vocab list to storage
func (r VocabRepo) Save(vl vocab.List) error {
	data, err := json.Marshal(vl.Words)
	if err != nil {
		return err
	}

	return os.WriteFile(r.Path, data, os.ModePerm)
}

// DefaultVocabFilepath returns the default filepath for where
// the vocab list may be stored on the filesystem
func DefaultVocabFilepath() string {
	return filepath.Join(defaultConfigDir(), vocabFile)
}
