package cmd

import (
	"errors"
	"fmt"
	"io"
	"math/rand"

	"github.com/caproven/termdict/dictionary"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type defineOptions struct {
	word   string
	limit  int
	random bool
}

// NewDefineCommand constructs the define command
func NewDefineCommand(cfg *Config) *cobra.Command {
	o := &defineOptions{}

	cmd := &cobra.Command{
		Use:   "define word | --random",
		Short: "Define a word",
		Long: `Lookup the definition for a given word.

Sample usage:
  termdict define organic
  termdict define --random`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.random && len(args) > 0 {
				return errors.New("can't use --random and a specified word")
			}

			if len(args) > 0 {
				o.word = args[0]
			}

			return o.run(cfg.Out, cfg.Vocab, cfg.Dict)
		},
	}

	cmd.Flags().IntVar(&o.limit, "limit", 0, "limit the number of entries to display")
	cmd.Flags().BoolVar(&o.random, "random", false, "define a random word from your vocab list")

	return cmd
}

func (o *defineOptions) run(out io.Writer, v VocabRepo, d Definer) error {
	word := o.word
	if o.random {
		var err error
		word, err = selectRandomWord(v)
		if err != nil {
			return err
		}
	}

	defs, err := d.Define(word)
	if err != nil {
		return err
	}

	printDefinition(out, word, defs, o.limit)
	return nil
}

func selectRandomWord(v VocabRepo) (string, error) {
	vl, err := v.Load()
	if err != nil {
		return "", err
	}

	if len(vl.Words) == 0 {
		return "", errors.New("no words in vocab list")
	}

	word := vl.Words[rand.Intn(len(vl.Words))]
	return word, nil
}

// printDefinition neatly prints a word along with its definitions. Allows limiting
// of definitions printed if limit > 0
func printDefinition(w io.Writer, word string, defs []dictionary.Definition, limit int) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Fprintln(w, green(word))

	blue := color.New(color.FgCyan).SprintFunc()
	if limit <= 0 {
		limit = len(defs)
	}
	for i, def := range defs {
		if i >= limit {
			break
		}
		fmt.Fprintf(w, "[%s] %s\n", blue(def.PartOfSpeech), def.Meaning)
	}
}
