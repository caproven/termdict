package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/rand"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type defineOptions struct {
	word       string
	random     bool
	randomSeed uint64
	save       bool
	output     string
	printers   map[string]defPrinter
}

// NewDefineCommand constructs the define command
func NewDefineCommand(cfg *Config) *cobra.Command {
	o := &defineOptions{
		printers: map[string]defPrinter{},
	}

	cmd := &cobra.Command{
		Use:   "define word | --random",
		Short: "Define a word",
		Long: `Lookup the definition for a given word.

Sample usage:
  termdict define organic
  termdict define --random`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.random {
				if len(args) > 0 {
					return errors.New("can't use --random and a specified word")
				}
			} else {
				if len(args) == 0 {
					return errors.New("must specify word")
				}
				o.word = args[0]
			}

			return o.run(cmd.Context(), cfg.Out, cfg.Vocab, cfg.Dict)
		},
	}

	o.registerPrinter(new(textPrinter), cmd)
	o.registerPrinter(new(jsonPrinter), cmd)

	cmd.Flags().BoolVar(&o.random, "random", false, "define a random word from your vocab list")
	cmd.Flags().Uint64Var(&o.randomSeed, "seed", 0, "rng seed making usage of --random deterministic")
	cmd.Flags().BoolVar(&o.save, "save", false, "add to the vocab list if the word can be defined")
	cmd.Flags().StringVarP(&o.output, "output", "o", "text", "output format; one of text, json")
	// Avoid attempting to save words already in the list.
	cmd.MarkFlagsMutuallyExclusive("save", "random")

	return cmd
}

func (o *defineOptions) run(ctx context.Context, out io.Writer, v VocabRepo, d Definer) error {
	// Don't call anything else if output format is invalid
	printer, err := o.getPrinter(o.output)
	if err != nil {
		return err
	}

	word := o.word
	if o.random {
		source := randSource(rand.Default{})
		if o.randomSeed != 0 {
			source = rand.NewRand(o.randomSeed)
		}

		var err error
		word, err = selectRandomWord(ctx, v, source)
		if err != nil {
			return err
		}
	}

	defs, err := d.Define(ctx, word)
	if err != nil {
		return err
	}

	if o.save {
		added, err := v.AddWordsToList(ctx, []string{word})
		if err != nil {
			return fmt.Errorf("save words to list: %w", err)
		}
		if len(added) > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "Saved word(s) to list: %q\n", added)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "Word already present in list, not saving")
		}
	}

	return printer.Print(out, word, defs)
}

func (o *defineOptions) registerPrinter(p defPrinter, cmd *cobra.Command) {
	o.printers[p.OutputType()] = p
}

func (o *defineOptions) getPrinter(output string) (defPrinter, error) {
	printer, ok := o.printers[strings.ToLower(output)]
	if !ok {
		return nil, fmt.Errorf("no printer registered for output %s", output)
	}
	return printer, nil
}

func selectRandomWord(ctx context.Context, v VocabRepo, randSource randSource) (string, error) {
	list, err := v.GetWordsInList(ctx)
	if err != nil {
		return "", fmt.Errorf("list words: %w", err)
	}

	if len(list) == 0 {
		return "", errors.New("no words found")
	}

	return list[randSource.IntN(len(list))], nil
}

type randSource interface {
	// IntN returns a random number in the interval [0, n).
	IntN(n int) int
}

type defPrinter interface {
	OutputType() string
	Print(w io.Writer, word string, defs []dictionary.Definition) error
}

type textPrinter struct {
}

func (p *textPrinter) OutputType() string {
	return "text"
}

func (p *textPrinter) Print(w io.Writer, word string, defs []dictionary.Definition) error {
	green := color.New(color.FgGreen).SprintFunc()
	if _, err := fmt.Fprintln(w, green(word)); err != nil {
		return err
	}

	blue := color.New(color.FgCyan).SprintFunc()
	for _, def := range defs {
		if _, err := fmt.Fprintf(w, "[%s] %s\n", blue(def.PartOfSpeech), def.Meaning); err != nil {
			return err
		}
	}

	return nil
}

type jsonPrinter struct{}

func (p *jsonPrinter) AddFlags(_ *cobra.Command) {
}

func (p *jsonPrinter) OutputType() string {
	return "json"
}

func (p *jsonPrinter) Print(w io.Writer, word string, defs []dictionary.Definition) error {
	composite := struct {
		Word        string
		Definitions []dictionary.Definition
	}{
		Word:        word,
		Definitions: defs,
	}

	data, err := json.MarshalIndent(composite, "", "\t")
	if err != nil {
		return err
	}
	if _, err = w.Write(data); err != nil {
		return err
	}
	if _, err = w.Write([]byte{'\n'}); err != nil {
		return err
	}

	return nil
}
