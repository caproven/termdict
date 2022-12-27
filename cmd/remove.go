package cmd

import (
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
  termdict remove efficacy
  termdict remove elegy chide`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.words = args

			return o.run(cfg.Out, cfg.Vocab)
		},
	}
	return cmd
}

func (o *removeOptions) run(out io.Writer, v VocabRepo) error {
	vl, err := v.Load()
	if err != nil {
		return err
	}

	for _, word := range o.words {
		if err := vl.RemoveWord(word); err != nil {
			return err
		}
	}

	return v.Save(vl)
}
