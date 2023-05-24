package cmd

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/options"
	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/ovhcloud"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/spf13/cobra"
)

// CreateCmd holds the "create" command flags.
type CreateCmd struct{}

// NewCreateCmd creates the "create" command.
func NewCreateCmd() *cobra.Command {
	cmd := &CreateCmd{}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return createCmd
}

// Run runs the "create" command logic.
func (cmd *CreateCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client, err := ovhcloud.NewClient(options.Authentication)
	if err != nil {
		return fmt.Errorf("can't create ovhcloud client: %w", err)
	}

	publicKeyBase, err := ssh.GetPublicKeyBase(options.MachineFolder)
	if err != nil {
		return fmt.Errorf("can't get public key: %w", err)
	}

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return fmt.Errorf("can't b64decode public key: %w", err)
	}

	opts := ovhcloud.CreateInstanceOptions{
		Name:      options.MachineID,
		Flavor:    options.Flavor,
		Image:     "Debian 10 - Docker",
		PublicKey: publicKey,
	}

	log.Infof("Starting machine creation with flavor %s", opts.Flavor)

	return client.CreateInstance(ctx, opts)
}
