package vocab

import (
	"fmt"
	"strings"
)

// List is vocab list, containing words
type List struct {
	Words []string
}

// AddWord adds a word to the vocab list. Duplicates are not
// permitted, and capitalization is ignored
func (vl *List) AddWord(w string) error {
	w = strings.ToLower(w)

	if vl.WordExists(w) {
		return fmt.Errorf("word '%s' already exists", w)
	}

	vl.Words = append(vl.Words, w)
	return nil
}

// RemoveWord removes a word from the vocab list. Capitalization
// is ignored
func (vl *List) RemoveWord(w string) error {
	w = strings.ToLower(w)

	idx := vl.idxOf(w)
	if idx == -1 {
		return fmt.Errorf("word '%s' doesn't exist", w)
	}
	vl.Words = append(vl.Words[:idx], vl.Words[idx+1:]...)

	return nil
}

// WordExists checks if a word exists in the vocab list
func (vl *List) WordExists(w string) bool {
	return vl.idxOf(w) != -1
}

// idxOf finds the index of the given word in the VocabList,
// or -1 if the word is not found
func (vl *List) idxOf(w string) int {
	for i, word := range vl.Words {
		if w == word {
			return i
		}
	}
	return -1
}
