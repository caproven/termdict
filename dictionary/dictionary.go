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

// PrintDefinition neatly prints a word along with its definitions. Allows limiting
// of definitions printed if limit > 0
func PrintDefinition(w io.Writer, word string, defs []Definition, limit int) {
	fmt.Fprintln(w, word)
	if limit <= 0 {
		limit = len(defs)
	}
	for i, def := range defs {
		if i >= limit {
			break
		}
		fmt.Fprintln(w, def)
	}
}
