package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/caproven/termdict/vocab"
	"github.com/spf13/cobra"
)

type importOptions struct{}

func NewImportCommand(cfg *Config) *cobra.Command {
	o := &importOptions{}

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import vocab events previously exported from another machine",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return o.run(cmd.Context(), cfg.Out, cfg.Vocab)
		},
	}

	return cmd
}

func (o *importOptions) run(ctx context.Context, out io.Writer, v VocabRepo) error {
	var events []vocab.Event
	dec := json.NewDecoder(os.Stdin)
	for {
		var event vocab.Event
		if err := dec.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("decode event: %w", err)
		}
		events = append(events, event)
	}

	if err := v.AddEvents(ctx, events); err != nil {
		return fmt.Errorf("import events: %w", err)
	}

	_, _ = fmt.Fprintln(out, "import complete")
	return nil
}
