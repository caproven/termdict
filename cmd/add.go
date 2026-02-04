package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type addOptions struct {
	words   []string
	noCheck bool
}

// NewAddCommand constructs the add command
func NewAddCommand(cfg *Config) *cobra.Command {
	o := &addOptions{}

	cmd := &cobra.Command{
		Use:   "add word ...",
		Short: "Add words to your vocab list",
		Long: `Add words to your personal vocab list.

Sample usage:
  termdict list add comeuppance
  termdict list add ameliorate entropy
  termdict list add omg --no-check`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.words = args

			return o.run(cmd.Context(), cfg.Out, cfg.Vocab, cfg.Dict)
		},
	}

	cmd.Flags().BoolVarP(&o.noCheck, "no-check", "n", false, "don't check that words can be defined before adding")

	return cmd
}

func (o *addOptions) run(ctx context.Context, out io.Writer, v VocabRepo, d Definer) error {
	if !o.noCheck {
		for _, word := range o.words {
			if _, err := d.Define(ctx, word); err != nil {
				return fmt.Errorf("define word %q to be added to list: %w", word, err)
			}
		}
	}

	added, err := v.AddWordsToList(ctx, o.words)
	if err != nil {
		return fmt.Errorf("add words to list: %w", err)
	}
	for _, word := range added {
		_, _ = fmt.Fprintf(out, "Added word %q\n", word)
	}

	return nil
}
