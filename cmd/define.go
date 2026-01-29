package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/caproven/termdict/dictionary"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type defineOptions struct {
	word     string
	random   bool
	output   string
	printers map[string]defPrinter
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
			if o.random && len(args) > 0 {
				return errors.New("can't use --random and a specified word")
			}

			if len(args) > 0 {
				o.word = args[0]
			}

			return o.run(cmd.Context(), cfg.Out, cfg.Vocab, cfg.Dict)
		},
	}

	o.registerPrinter(new(textPrinter), cmd)
	o.registerPrinter(new(jsonPrinter), cmd)

	cmd.Flags().BoolVar(&o.random, "random", false, "define a random word from your vocab list")
	cmd.Flags().StringVarP(&o.output, "output", "o", "text", "output format; one of text, json")

	return cmd
}

func (o *defineOptions) run(ctx context.Context, out io.Writer, v VocabRepo, d Definer) error {
	word := o.word
	if o.random {
		var err error
		word, err = selectRandomWord(ctx, v)
		if err != nil {
			return err
		}
	}

	defs, err := d.Define(ctx, word)
	if err != nil {
		return err
	}

	printer, err := o.getPrinter(o.output)
	if err != nil {
		return err
	}
	return printer.Print(out, word, defs)
}

func (o *defineOptions) registerPrinter(p defPrinter, cmd *cobra.Command) {
	o.printers[p.OutputType()] = p
	p.AddFlags(cmd)
}

func (o *defineOptions) getPrinter(output string) (defPrinter, error) {
	printer, ok := o.printers[strings.ToLower(output)]
	if !ok {
		return nil, fmt.Errorf("no printer registered for output %s", output)
	}
	return printer, nil
}

func selectRandomWord(ctx context.Context, v VocabRepo) (string, error) {
	list, err := v.GetWordsInList(ctx)
	if err != nil {
		return "", fmt.Errorf("list words: %w", err)
	}

	if len(list) == 0 {
		return "", errors.New("no words found")
	}

	word := list[rand.Intn(len(list))]
	return word, nil
}

type defPrinter interface {
	AddFlags(*cobra.Command)
	OutputType() string
	Print(w io.Writer, word string, defs []dictionary.Definition) error
}

type textPrinter struct {
	limit int
}

func (p *textPrinter) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&p.limit, "limit", 0, "limit the number of entries to display in text format")
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
	limit := p.limit
	if limit <= 0 || limit > len(defs) {
		limit = len(defs)
	}
	for i, def := range defs[:limit] {
		if i >= limit {
			break
		}
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
