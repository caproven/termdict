package cmd

import (
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/spf13/cobra"
)

type defineOptions struct {
	word string
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
	return cmd
}

func (o *defineOptions) run(out io.Writer, d dictionary.API) error {
	defs, err := d.Define(o.word)
	if err != nil {
		return err
	}

	dictionary.PrintDefinition(out, o.word, defs)

	return nil
}
