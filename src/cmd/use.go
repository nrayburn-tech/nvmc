package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"nvmc/util"
	"os"
)

type useCmd struct {
	command *cobra.Command
}

func newUseCmd() *useCmd {
	cmd := &useCmd{}
	cmd.command = &cobra.Command{
		Use:   "use <version>",
		Short: "Set <version> to the current node version.",
		Example: `# Use version 18.2.0.
$ nvmc use 18.2.0`,
		Args: cobra.ExactArgs(1),
		RunE: cmd.run(),
	}

	return cmd
}

func (c *useCmd) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return use(version)
	}
}

func use(version string) error {
	version, err := util.NormalizeVersion(version)
	if err != nil {
		return err
	}

	currentVersionDir, err := util.GetVersionPath(version)
	if err != nil {
		return err
	}
	currentVersionExeDir, err := util.GetVersionExePath(version)
	if err != nil {
		return err
	}
	nodeSymLink, err := util.GetSymLinkPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(currentVersionDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("requested installation " + version + " does not exist")
		}
		return err
	}

	if err := os.Remove(nodeSymLink); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Symlink(currentVersionExeDir, nodeSymLink); err != nil {
		return err
	}

	return nil
}
