package dictionary

import (
	"fmt"
	"io"
)

type Entry struct {
	PartOfSpeech string
	Definition   string
}

func (e Entry) String() string {
	return fmt.Sprintf("[%s] %s", e.PartOfSpeech, e.Definition)
}

type Definer interface {
	Define(word string) ([]Entry, error)
}

func Default() Definer {
	return NewDictionaryAPI()
}

func PrintDefinition(w io.Writer, word string, defs []Entry) {
	fmt.Fprintln(w, word)
	for _, def := range defs {
		fmt.Fprintln(w, def)
	}
}
