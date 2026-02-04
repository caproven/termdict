package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type removeOptions struct {
	words []string
}

// NewRemoveCommand constructs the remove command
func NewRemoveCommand(cfg *Config) *cobra.Command {
	o := &removeOptions{}

	cmd := &cobra.Command{
		Use:   "remove word ...",
		Short: "Remove words from your vocab list",
		Long: `Remove words from your personal vocab list.

Sample usage:
  termdict list remove efficacy
  termdict list remove elegy chide`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.words = args

			return o.run(cmd.Context(), cfg.Out, cfg.Vocab)
		},
	}
	return cmd
}

func (o *removeOptions) run(ctx context.Context, out io.Writer, v VocabRepo) error {
	removed, err := v.RemoveWordsFromList(ctx, o.words)
	if err != nil {
		return fmt.Errorf("remove words from list: %w", err)
	}
	for _, word := range removed {
		_, _ = fmt.Fprintf(out, "Removed word %q\n", word)
	}
	return nil
}
