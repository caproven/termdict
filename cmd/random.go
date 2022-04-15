package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
	"github.com/spf13/cobra"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type randomOptions struct {
}

func newRandomCommand(cfg *Config) *cobra.Command {
	o := &randomOptions{}

	cmd := &cobra.Command{
		Use:   "random",
		Short: "Define a random word from your vocab list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(cfg.Out, cfg.Storage)
		},
	}
	return cmd
}

func (o *randomOptions) run(out io.Writer, s vocab.Storage) error {
	vl, err := s.Load()
	if err != nil {
		return nil
	}

	if len(vl.Words) == 0 {
		fmt.Fprintln(out, "no words in vocab list")
		return nil
	}

	word := vl.Words[rand.Intn(len(vl.Words))]
	dict := dictionary.Default()

	defs, err := dict.Define(word)
	if err != nil {
		return err
	}

	dictionary.PrintDefinition(out, word, defs)

	return nil
}
