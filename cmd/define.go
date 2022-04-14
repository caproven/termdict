package cmd

import (
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/spf13/cobra"
)

type defineOptions struct {
	word string
}

func newDefineCommand(cfg *Config) *cobra.Command {
	o := &defineOptions{}

	cmd := &cobra.Command{
		Use:   "define word",
		Short: "Define a word",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.word = args[0]

			return o.run(cfg.Out)
		},
	}
	return cmd
}

func (o *defineOptions) run(out io.Writer) error {
	dict := dictionary.Default()

	defs, err := dict.Define(o.word)
	if err != nil {
		return err
	}

	dictionary.PrintDefinition(out, o.word, defs)

	return nil
}
