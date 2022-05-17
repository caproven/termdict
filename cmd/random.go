package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage"
	"github.com/spf13/cobra"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type randomOptions struct {
}

// NewRandomCommand constructs the random command
func NewRandomCommand(cfg *Config) *cobra.Command {
	o := &randomOptions{}

	cmd := &cobra.Command{
		Use:   "random",
		Short: "Define a random word from your vocab list",
		Long: `Define a word at random from your personal vocab list.

Sample usage:
  termdict random`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(cfg.Out, cfg.Vocab, cfg.Dict)
		},
	}
	return cmd
}

func (o *randomOptions) run(out io.Writer, v storage.VocabRepo, d dictionary.API) error {
	vl, err := v.Load()
	if err != nil {
		return nil
	}

	if len(vl.Words) == 0 {
		fmt.Fprintln(out, "no words in vocab list")
		return nil
	}

	word := vl.Words[rand.Intn(len(vl.Words))]

	defs, err := d.Define(word)
	if err != nil {
		return err
	}

	dictionary.PrintDefinition(out, word, defs)

	return nil
}
