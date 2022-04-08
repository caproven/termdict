package internal

import (
	"errors"
	"strings"
)

type VocabList []string

// TODO add alt funcs for multiple words so you can add/remove
// without repeatedly opening file

func AddWord(w string) error {
	words, err := ListWords()
	if err != nil {
		return err
	}

	w = strings.ToLower(w)

	for _, word := range words {
		if w == word {
			return errors.New("cannot add duplicate word")
		}
	}

	words = append(words, w)

	return writeVocabFile(words)
}

func RemoveWord(w string) error {
	words, err := ListWords()
	if err != nil {
		return err
	}

	w = strings.ToLower(w)

	idx := -1
	for i, word := range words {
		if w == word {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("cannot delete word not in vocab list")
	}

	words = append(words[:idx], words[idx+1:]...)

	return writeVocabFile(words)
}

func ListWords() (VocabList, error) {
	v, err := readVocabFile()
	if err != nil {
		return nil, err
	}
	return v, nil
}
