package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type listOptions struct {
}

// NewListCommand constructs the list command
func NewListCommand(cfg *Config) *cobra.Command {
	o := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the words in your vocab list",
		Long: `List the words in your personal vocab list.

Sample usage:
  termdict list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(cmd.Context(), cfg.Out, cfg.Vocab)
		},
	}

	cmd.AddCommand(NewAddCommand(cfg))
	cmd.AddCommand(NewRemoveCommand(cfg))

	return cmd
}

// TODO support json
func (o *listOptions) run(ctx context.Context, out io.Writer, v VocabRepo) error {
	words, err := v.GetWordsInList(ctx)
	if err != nil {
		return fmt.Errorf("list words: %w", err)
	}

	if len(words) == 0 {
		fmt.Fprintln(out, "no words in vocab list")
		return nil
	}

	for _, word := range words {
		fmt.Fprintln(out, word)
	}

	return nil
}
