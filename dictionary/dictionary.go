package dictionary

import (
	"fmt"
	"io"
)

// Entry is a single dictionary entry for a word
type Entry struct {
	PartOfSpeech string
	Definition   string
}

// String formats an Entry as a string
func (e Entry) String() string {
	return fmt.Sprintf("[%s] %s", e.PartOfSpeech, e.Definition)
}

// PrintDefinition neatly prints a word along with its definitions
func PrintDefinition(w io.Writer, word string, defs []Entry) {
	fmt.Fprintln(w, word)
	for _, def := range defs {
		fmt.Fprintln(w, def)
	}
}
