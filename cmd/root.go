package cmd

import (
	"io"

	"github.com/caproven/termdict/storage"
	"github.com/spf13/cobra"
)

type Config struct {
	Out     io.Writer
	Storage storage.VocabStorage
}

func NewRootCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "termdict",
		Short: "A small dictionary tool for the command line",
	}

	cmd.AddCommand(newAddCmd(cfg))
	cmd.AddCommand(newDefineCommand(cfg))
	cmd.AddCommand(newListCommand(cfg))
	cmd.AddCommand(newRandomCommand(cfg))
	cmd.AddCommand(newRemoveCommand(cfg))

	return cmd
}
