package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type exportOptions struct{}

func NewExportCommand(cfg *Config) *cobra.Command {
	o := &exportOptions{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export vocab events for import on another machine",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return o.run(cmd.Context(), cfg.Out, cfg.Vocab)
		},
	}

	return cmd
}

func (o *exportOptions) run(ctx context.Context, out io.Writer, v VocabRepo) error {
	events, err := v.GetEvents(ctx)
	if err != nil {
		return fmt.Errorf("export events: %w", err)
	}

	enc := json.NewEncoder(out)
	for _, event := range events {
		if err := enc.Encode(event); err != nil {
			return fmt.Errorf("encode event: %w", err)
		}
	}

	return nil
}
