package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/caproven/termdict/vocab"
)

const configSubdir string = "termdict"
const vocabFile string = "vocab.json"

type VocabStorage interface {
	Read() (vocab.VocabList, error)
	Write(vocab.VocabList) error
}

type VocabFile struct {
	Path string
}

func (vf VocabFile) Read() (vocab.VocabList, error) {
	var vl vocab.VocabList

	if _, err := os.Stat(vf.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// file doesn't exist, just return empty list
			return vl, nil
		}
		return vocab.VocabList{}, err
	}

	data, err := os.ReadFile(vf.Path)
	if err != nil {
		return vocab.VocabList{}, err
	}

	if err := json.Unmarshal(data, &vl.Words); err != nil {
		return vocab.VocabList{}, err
	}

	return vl, nil
}

func (vf VocabFile) Write(vl vocab.VocabList) error {
	data, err := json.Marshal(vl.Words)
	if err != nil {
		return err
	}

	if err := createConfigDir(); err != nil {
		return err
	}

	return os.WriteFile(vf.Path, data, os.ModePerm)
}

func DefaultVocabFile() string {
	return filepath.Join(defaultConfigDir(), vocabFile)
}

func defaultConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, configSubdir)
}

func createConfigDir() error {
	return os.MkdirAll(defaultConfigDir(), os.ModePerm)
}
