package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"nvmc/util"
	"os"
)

type uninstallCmd struct {
	command       *cobra.Command
	globalOpts    globalOpts
	uninstallOpts uninstallOpts
}

func newUninstallCmd(globalOpts globalOpts) *uninstallCmd {
	cmd := &uninstallCmd{}
	cmd.command = &cobra.Command{
		Use:   "uninstall <version>",
		Short: "Uninstall <version>.",
		Example: `# Uninstall version 18.2.0.
$ nvmc uninstall 18.2.0`,
		Args: cobra.ExactArgs(1),
		RunE: cmd.run(),
	}

	cmd.globalOpts = globalOpts

	return cmd
}

func (c *uninstallCmd) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return uninstall(version, c.globalOpts, c.uninstallOpts)
	}
}

func uninstall(version string, globalOpts globalOpts, uninstallOpts uninstallOpts) error {
	version, err := util.NormalizeVersion(version)
	if err != nil {
		return err
	}

	currentVersionDir, err := util.GetVersionPath(version)
	if err != nil {
		return err
	}

	stats, err := os.Stat(currentVersionDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return errors.New("Version does not exist. Path: " + currentVersionDir)
	} else if err != nil {
		return err
	}
	if !stats.IsDir() {
		return errors.New("Version path already exists and is not a directory. Path: " + currentVersionDir)
	}

	if err := os.RemoveAll(currentVersionDir); err != nil {
		return err
	}

	return nil
}
