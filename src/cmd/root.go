package cmd

import (
	"github.com/spf13/cobra"
	"nvmc/util"
	"os"
)

type rootCmd struct {
	command    *cobra.Command
	globalOpts globalOpts
}

func newRootCmd() *rootCmd {
	cmd := &rootCmd{}
	cmd.command = &cobra.Command{
		Use:     "nvmc",
		Short:   "Install and manage multiple versions of node",
		Version: util.VERSION,
	}

	cmd.command.Flags().StringVar(&cmd.globalOpts.downloadUrl, "download-url", defaultGlobalOpts.downloadUrl, "Specify a custom base URL.")
	cmd.command.Flags().BoolVar(&cmd.globalOpts.followRedirects, "follow-redirects", defaultGlobalOpts.followRedirects, "Follow redirects when downloading files.")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := newRootCmd()
	// Hide the completions command, but keep it available
	rootCmd.command.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.command.AddCommand(newInstallCmd(rootCmd.globalOpts).command)
	rootCmd.command.AddCommand(newListCmd(rootCmd.globalOpts).command)
	rootCmd.command.AddCommand(newUninstallCmd(rootCmd.globalOpts).command)
	rootCmd.command.AddCommand(newUseCmd(rootCmd.globalOpts).command)

	err := rootCmd.command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
