package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nvmc",
	Short: "Install and manage multiple versions of node",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Hide the completions command, but keep it available
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(newInstallCmd().command)
	rootCmd.AddCommand(newListCmd().command)
	rootCmd.AddCommand(newUninstallCmd().command)
	rootCmd.AddCommand(newUseCmd().command)
}
