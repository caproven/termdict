package cmd

import (
	"io"

	"github.com/caproven/termdict/storage"
	"github.com/spf13/cobra"
)

type removeOptions struct {
	words []string
}

func newRemoveCommand(cfg *Config) *cobra.Command {
	o := &removeOptions{}

	cmd := &cobra.Command{
		Use:   "remove word ...",
		Short: "Remove words from your vocab list",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.words = args

			return o.run(cfg.Out, cfg.Storage)
		},
	}
	return cmd
}

func (o *removeOptions) run(out io.Writer, s storage.VocabStorage) error {
	vl, err := s.Read()
	if err != nil {
		return err
	}

	for _, word := range o.words {
		if err := vl.RemoveWord(word); err != nil {
			return err
		}
	}

	if err := s.Write(vl); err != nil {
		return err
	}

	return nil
}
