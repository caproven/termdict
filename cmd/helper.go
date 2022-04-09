package cmd

import (
	"fmt"
	"io"

	"github.com/caproven/termdict/internal/dictionary"
)

func printDefinition(w io.Writer, word string, defs []dictionary.Definition) {
	fmt.Fprintln(w, word)
	for _, def := range defs {
		fmt.Fprintf(w, "[%s] %s\n", def.PartOfSpeech, def.Meaning)
	}
}
