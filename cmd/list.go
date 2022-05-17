package cmd

import (
	"fmt"
	"io"

	"github.com/caproven/termdict/storage"
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
			return o.run(cfg.Out, cfg.Vocab)
		},
	}
	return cmd
}

func (o *listOptions) run(out io.Writer, v storage.VocabRepo) error {
	vl, err := v.Load()
	if err != nil {
		return err
	}

	if len(vl.Words) == 0 {
		fmt.Fprintln(out, "no words in vocab list")
		return nil
	}

	for _, word := range vl.Words {
		fmt.Fprintln(out, word)
	}

	return nil
}
