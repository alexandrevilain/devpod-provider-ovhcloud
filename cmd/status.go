package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/options"
	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/ovhcloud"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/spf13/cobra"
)

// StatusCmd holds the cmd flags
type StatusCmd struct{}

// NewStatusCmd defines a command
func NewStatusCmd() *cobra.Command {
	cmd := &StatusCmd{}
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Retrieve the status of an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return statusCmd
}

// Run runs the command logic
func (cmd *StatusCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client, err := ovhcloud.NewClient(options.Authentication)
	if err != nil {
		return err
	}

	status, err := client.GetInstanceStatus(ctx, options.MachineID)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, OVHCloudStatusToDevPodStatus(status))
	return err
}

func OVHCloudStatusToDevPodStatus(status ovhcloud.InstanceStatus) string {
	switch status {
	case ovhcloud.InstanceActive:
		return "Running"
	case ovhcloud.InstanceStopped, ovhcloud.InstanceShutOff:
		return "Stopped"
	default:
		return "Busy"
	}
}
