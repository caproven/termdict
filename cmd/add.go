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
	vl, err := v.Load()
	if err != nil {
		return err
	}

	for _, word := range o.words {
		if !o.noCheck {
			if _, err := d.Define(ctx, word); err != nil {
				return fmt.Errorf("failed to add word '%s'; couldn't find a definition", word)
			}
		}

		if err := vl.AddWord(word); err != nil {
			return err
		}
	}

	return v.Save(vl)
}
