package cmd

import (
	"context"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/options"
	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/ovhcloud"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/spf13/cobra"
)

// DeleteCmd holds the "delete" command flags.
type DeleteCmd struct{}

// NewDeleteCmd creates the "delete" command.
func NewDeleteCmd() *cobra.Command {
	cmd := &DeleteCmd{}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return deleteCmd
}

// Run runs the "delete" command logic.
func (cmd *DeleteCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client, err := ovhcloud.NewClient(options.Authentication)
	if err != nil {
		return err
	}

	return client.DeleteInstance(ctx, options.MachineID)
}
