package vocab

import (
	"fmt"
	"strings"
)

type VocabList struct {
	Words []string
}

func (vl *VocabList) AddWord(w string) error {
	w = strings.ToLower(w)

	if vl.wordExists(w) {
		return fmt.Errorf("word '%s' already exists", w)
	}

	vl.Words = append(vl.Words, w)
	return nil
}

func (vl *VocabList) RemoveWord(w string) error {
	w = strings.ToLower(w)

	idx := vl.idxOf(w)
	if idx == -1 {
		return fmt.Errorf("word '%s' doesn't exist", w)
	}
	vl.Words = append(vl.Words[:idx], vl.Words[idx+1:]...)

	return nil
}

// idxOf finds the index of the given word in the VocabList,
// or -1 if the word is not found
func (vl *VocabList) idxOf(w string) int {
	for i, word := range vl.Words {
		if w == word {
			return i
		}
	}
	return -1
}

func (vl *VocabList) wordExists(w string) bool {
	return vl.idxOf(w) != -1
}
