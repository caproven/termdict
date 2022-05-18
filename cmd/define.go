package cmd

import (
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/storage"
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

			return o.run(cfg.Out, cfg.Cache, cfg.Dict)
		},
	}
	return cmd
}

func (o *defineOptions) run(out io.Writer, c storage.CacheRepo, d dictionary.API) error {
	cache, err := c.Load()
	if err != nil {
		return err
	}

	defs, ok := cache[o.word]
	if !ok {
		defs, err = d.Define(o.word)
		if err != nil {
			return err
		}

		cache[o.word] = defs
		if err := c.Save(cache); err != nil {
			return err
		}
	}

	dictionary.PrintDefinition(out, o.word, defs)

	return nil
}
