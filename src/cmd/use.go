package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"nvmc/util"
	"os"
	"path/filepath"
	"runtime"
)

type useCmd struct {
	command    *cobra.Command
	globalOpts globalOpts
	useOpts    useOpts
}

func newUseCmd(globalOpts globalOpts) *useCmd {
	cmd := &useCmd{}
	cmd.command = &cobra.Command{
		Use:   "use <version>",
		Short: "Set <version> to the current node version.",
		Example: `# Use version 18.2.0.
$ nvmc use 18.2.0`,
		Args: cobra.ExactArgs(1),
		RunE: cmd.run(),
	}

	cmd.globalOpts = globalOpts

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

	// TODO: Validate the version
	currentVersionDir, err := util.GetVersionPath(version)
	if err != nil {
		return err
	}

	currentVersionStats, err := os.Stat(currentVersionDir)
	if err != nil {
		return err
	}
	if !currentVersionStats.IsDir() {
		return errors.New("Version path already exists and is not a directory. Path: " + currentVersionDir)
	}

	nodeSymLink, err := util.GetSymLinkPath()
	if err != nil {
		return err
	}

	if err := os.Remove(nodeSymLink); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	installationInfo, err := util.GetInstallationInfo(version)
	if err != nil {
		return err
	}

	var symLinkTarget string
	if runtime.GOOS == "windows" {
		symLinkTarget = filepath.Join(currentVersionDir, installationInfo.FileNameWithoutExtension)
	} else {
		symLinkTarget = filepath.Join(currentVersionDir, installationInfo.FileNameWithoutExtension, "bin")
	}
	if err := os.Symlink(symLinkTarget, nodeSymLink); err != nil {
		return err
	}

	fmt.Println("now using node " + version)

	return nil
}
