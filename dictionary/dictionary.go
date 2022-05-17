package dictionary

import (
	"fmt"
	"io"
)

// Definition is a single dictionary entry for a word
type Definition struct {
	PartOfSpeech string
	Meaning      string
}

// String formats a word definition as a string
func (def Definition) String() string {
	return fmt.Sprintf("[%s] %s", def.PartOfSpeech, def.Meaning)
}

// PrintDefinition neatly prints a word along with its definitions
func PrintDefinition(w io.Writer, word string, defs []Definition) {
	fmt.Fprintln(w, word)
	for _, def := range defs {
		fmt.Fprintln(w, def)
	}
}
