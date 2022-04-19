package cmd

import (
	"io"

	"github.com/caproven/termdict/dictionary"
	"github.com/caproven/termdict/vocab"
	"github.com/spf13/cobra"
)

// Config represents the CLI configuration
type Config struct {
	Out     io.Writer
	Storage vocab.Storage
	DictAPI dictionary.API
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
