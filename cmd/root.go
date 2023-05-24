package cmd

import (
	"os"
	"os/exec"

	log "github.com/loft-sh/devpod/pkg/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// NewRootCmd creates a new root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "devpod-provider-ovhcloud",
		Short:         "ovhcloud provider commands",
		SilenceErrors: true,
		SilenceUsage:  true,

		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			log.Default.MakeRaw()
			return nil
		},
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd := NewRootCmd()
	rootCmd.AddCommand(NewCreateCmd())
	rootCmd.AddCommand(NewStatusCmd())
	rootCmd.AddCommand(NewDeleteCmd())
	rootCmd.AddCommand(NewStartCmd())
	rootCmd.AddCommand(NewStopCmd())
	rootCmd.AddCommand(NewCommandCmd())
	rootCmd.AddCommand(NewInitCmd())

	// execute command
	err := rootCmd.Execute()
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			os.Exit(exitErr.ExitStatus())
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			if len(exitErr.Stderr) > 0 {
				log.Default.ErrorStreamOnly().Error(string(exitErr.Stderr))
			}
			os.Exit(exitErr.ExitCode())
		}

		log.Default.Fatal(err)
	}
}
