package cmd

import (
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
	"github.com/spf13/cobra"
)

// Config represents the CLI configuration
type Config struct {
	Out   io.Writer
	Vocab VocabRepo
	Cache Cache
	Dict  Definer
}

type Definer interface {
	Define(word string) ([]dictionary.Definition, error)
}

type Cache interface {
	Contains(word string) (bool, error)
	Save(word string, defs []dictionary.Definition) error
	Lookup(word string) ([]dictionary.Definition, error)
}

type VocabRepo interface {
	Load() (vocab.List, error)
	Save(vl vocab.List) error
}

// NewRootCmd creates and returns an instance of the root command
func NewRootCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "termdict",
		Short: "A small dictionary tool for the command line",
	}

	cmd.AddCommand(NewAddCommand(cfg))
	cmd.AddCommand(NewDefineCommand(cfg))
	cmd.AddCommand(NewListCommand(cfg))
	cmd.AddCommand(NewRandomCommand(cfg))
	cmd.AddCommand(NewRemoveCommand(cfg))

	return cmd
}
