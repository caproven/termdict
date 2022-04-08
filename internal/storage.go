package internal

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

const configSubdir string = "termdict"
const vocabFile string = "vocab.json"

func defaultConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, configSubdir)
}

func defaultVocabFile() string {
	return filepath.Join(defaultConfigDir(), vocabFile)
}

func createConfigDir() error {
	return os.MkdirAll(defaultConfigDir(), os.ModePerm)
}

func readVocabFile() (VocabList, error) {
	var v VocabList

	fp := filepath.Join(defaultConfigDir(), vocabFile)
	if _, err := os.Stat(fp); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// file doesn't exist, just return empty list
			return v, nil
		}
		return nil, err
	}

	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func writeVocabFile(v VocabList) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if err := createConfigDir(); err != nil {
		return err
	}

	return os.WriteFile(defaultVocabFile(), data, os.ModePerm)
}
