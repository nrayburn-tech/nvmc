package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"nvmc/util"
	"os"
)

type uninstallCmd struct {
	command *cobra.Command
}

func newUninstallCmd() *uninstallCmd {
	cmd := &uninstallCmd{}
	cmd.command = &cobra.Command{
		Use:   "uninstall <version>",
		Short: "Uninstall <version>.",
		Example: `# Uninstall version 18.2.0.
$ nvmc uninstall 18.2.0`,
		Args: cobra.ExactArgs(1),
		RunE: cmd.run(),
	}

	return cmd
}

func (c *uninstallCmd) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return uninstall(version)
	}
}

func uninstall(version string) error {
	version, err := util.NormalizeVersion(version)
	if err != nil {
		return err
	}

	// TODO: Validate the version

	currentVersionDir, err := util.GetVersionPath(version)
	if err != nil {
		return err
	}

	if _, err := os.Stat(currentVersionDir); err == nil {
		if err := os.RemoveAll(currentVersionDir); err != nil {
			return err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return errors.New("requested installation " + version + " does not exit")
	}

	return nil
}
