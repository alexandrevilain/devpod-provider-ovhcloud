package cmd

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/ovhcloud"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/options"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// CommandCmd holds the "command" commands flags.
type CommandCmd struct{}

// NewCommandCmd creates the "command" command.
func NewCommandCmd() *cobra.Command {
	cmd := &CommandCmd{}
	commandCmd := &cobra.Command{
		Use:   "command",
		Short: "Run a command on the instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return commandCmd
}

// Run runs the "command" command logic.
func (cmd *CommandCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	command := os.Getenv("COMMAND")
	if command == "" {
		return fmt.Errorf("COMMAND environment variable is missing")
	}

	// get private key
	privateKey, err := ssh.GetPrivateKeyRawBase(options.MachineFolder)
	if err != nil {
		return fmt.Errorf("can't get private key: %w", err)
	}

	client, err := ovhcloud.NewClient(options.Authentication)
	if err != nil {
		return err
	}

	instance, err := client.GetInstanceByName(ctx, options.MachineID)
	if err != nil {
		return err
	}

	externalIP, err := ovhcloud.PublicIP(instance)
	if err != nil {
		return err
	}

	sshClient, err := ssh.NewSSHClient("debian", net.JoinHostPort(externalIP, "22"), privateKey)
	if err != nil {
		return errors.Wrap(err, "create ssh client")
	}
	defer sshClient.Close()

	return ssh.Run(context.Background(), sshClient, command, os.Stdin, os.Stdout, os.Stderr)
}
