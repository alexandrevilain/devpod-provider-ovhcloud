package cmd

import (
	"context"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/options"
	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/ovhcloud"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/spf13/cobra"
)

// StopCmd holds the "stop" command flags.
type StopCmd struct{}

// NewStopCmd creates the "stop" command.
func NewStopCmd() *cobra.Command {
	cmd := &StopCmd{}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return stopCmd
}

// Run runs "stop" the command logic.
func (cmd *StopCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client, err := ovhcloud.NewClient(options.Authentication)
	if err != nil {
		return err
	}

	return client.StopInstance(ctx, options.MachineID)
}
