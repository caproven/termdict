package cmd

import (
	"fmt"
	"io"

	"github.com/caproven/termdict/storage"
	"github.com/spf13/cobra"
)

type listOptions struct {
}

func newListCommand(cfg *Config) *cobra.Command {
	o := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the words in your vocab list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(cfg.Out, cfg.Storage)
		},
	}
	return cmd
}

func (o *listOptions) run(out io.Writer, s storage.VocabStorage) error {
	vl, err := s.Read()
	if err != nil {
		return err
	}

	if len(vl.Words) == 0 {
		fmt.Fprintln(out, "no words in vocab list")
		return nil
	}

	for _, word := range vl.Words {
		fmt.Fprintln(out, word)
	}

	return nil
}
