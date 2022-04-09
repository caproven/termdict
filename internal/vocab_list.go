package internal

import (
	"fmt"
	"strings"
)

type VocabList []string

func AddWords(words ...string) error {
	list, err := ListWords()
	if err != nil {
		return err
	}

	for _, word := range words {
		word = strings.ToLower(word)

		if exists, _ := wordExists(list, word); exists {
			return fmt.Errorf("word '%s' already exists", word)
		}
		list = append(list, words...)
	}

	return writeVocabFile(list)
}

func RemoveWords(words ...string) error {
	list, err := ListWords()
	if err != nil {
		return err
	}

	for _, word := range words {
		word = strings.ToLower(word)

		exists, idx := wordExists(list, word)
		if !exists {
			return fmt.Errorf("word '%s' doesn't exist", word)
		}
		list = append(list[:idx], list[idx+1:]...)
	}

	return writeVocabFile(words)
}

func ListWords() (VocabList, error) {
	v, err := readVocabFile()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func wordExists(list VocabList, w string) (bool, int) {
	for i, existingWord := range list {
		if w == existingWord {
			return true, i
		}
	}
	return false, 0
}
