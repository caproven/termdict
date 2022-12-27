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

// FileVocabRepo lets you fetch or save a vocab list
type FileVocabRepo struct {
	path string
}

// Load a vocab list from storage
func (r FileVocabRepo) Load() (vocab.List, error) {
	var vl vocab.List

	if _, err := os.Stat(r.path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// file doesn't exist, just return empty vocab.List
			return vl, nil
		}
		return vocab.List{}, err
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		return vocab.List{}, err
	}

	if err := json.Unmarshal(data, &vl.Words); err != nil {
		return vocab.List{}, err
	}

	return vl, nil
}

// Save a vocab list to storage
func (r FileVocabRepo) Save(vl vocab.List) error {
	data, err := json.Marshal(vl.Words)
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, data, os.ModePerm)
}

// defaultVocabFilepath returns the default filepath for where
// the vocab list may be stored on the filesystem
func defaultVocabFilepath() string {
	return filepath.Join(defaultConfigDir(), vocabFile)
}

func NewDefaultVocabRepo() FileVocabRepo {
	return FileVocabRepo{
		path: defaultVocabFilepath(),
	}
}
