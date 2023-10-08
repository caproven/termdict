package cmd

import (
	"fmt"
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type defineOptions struct {
	word  string
	limit int
}

// NewDefineCommand constructs the define command
func NewDefineCommand(cfg *Config) *cobra.Command {
	o := &defineOptions{}

	cmd := &cobra.Command{
		Use:   "define word",
		Short: "Define a word",
		Long: `Lookup the definition for a given word.

Sample usage:
  termdict define organic
  termdict define symphonic`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.word = args[0]

			return o.run(cfg.Out, cfg.Dict)
		},
	}

	cmd.Flags().IntVar(&o.limit, "limit", 0, "limit the number of entries to display")

	return cmd
}

func (o *defineOptions) run(out io.Writer, d Definer) error {
	defs, err := d.Define(o.word)
	if err != nil {
		return err
	}

	printDefinition(out, o.word, defs, o.limit)
	return nil
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
