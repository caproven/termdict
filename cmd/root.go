package cmd

import (
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Config represents the CLI configuration
type Config struct {
	Out   io.Writer
	Vocab VocabRepo
	Dict  Definer
}

type Definer interface {
	Define(word string) ([]dictionary.Definition, error)
}

type VocabRepo interface {
	Load() (vocab.List, error)
	Save(vl vocab.List) error
}

type rootOptions struct {
	noColor bool
}

// NewRootCmd creates and returns an instance of the root command
func NewRootCmd(cfg *Config) *cobra.Command {
	o := &rootOptions{}

	cmd := &cobra.Command{
		Use:   "termdict",
		Short: "A small dictionary tool for the command line",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if o.noColor {
				color.NoColor = true
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&o.noColor, "no-color", false, "disable colorized output")

	cmd.AddCommand(NewAddCommand(cfg))
	cmd.AddCommand(NewDefineCommand(cfg))
	cmd.AddCommand(NewListCommand(cfg))
	cmd.AddCommand(NewRemoveCommand(cfg))

	return cmd
}
