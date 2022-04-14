package cmd

import (
	"fmt"
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
	"github.com/spf13/cobra"
)

type addOptions struct {
	words   []string
	noCheck bool
}

func newAddCmd(cfg *Config) *cobra.Command {
	o := &addOptions{}

	cmd := &cobra.Command{
		Use:   "add word ...",
		Short: "Add words to your vocab list",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.words = args

			return o.run(cfg.Out, cfg.Storage)
		},
	}

	cmd.Flags().BoolVarP(&o.noCheck, "no-check", "n", false, "don't check that words can be defined before adding")

	return cmd
}

func (o *addOptions) run(out io.Writer, s vocab.Storage) error {
	vl, err := s.Read()
	if err != nil {
		return err
	}

	dict := dictionary.Default()

	for _, word := range o.words {
		if !o.noCheck {
			if _, err := dict.Define(word); err != nil {
				return fmt.Errorf("failed to add word '%s'; couldn't find a definition", word)
			}
		}

		if err := vl.AddWord(word); err != nil {
			return err
		}
	}

	if err := s.Write(vl); err != nil {
		return err
	}
	return nil
}
