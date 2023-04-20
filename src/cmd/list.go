package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"nvmc/util"
	"os"
	"strings"
)

type listCmd struct {
	command *cobra.Command
}

func newListCmd() *listCmd {
	cmd := &listCmd{}
	cmd.command = &cobra.Command{
		Use:     "list",
		Short:   "List all installed node versions.",
		Example: `$ nvmc list`,
		Args:    cobra.ExactArgs(0),
		RunE:    cmd.run(),
	}

	return cmd
}

func (c *listCmd) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return list()
	}
}

func list() error {
	versions, err := retrieveVersions()
	if err != nil {
		return err
	} else if len(versions) == 0 {
		return errors.New("no versions installed")
	}
	current, err := currentVersion()

	for _, version := range versions {
		fmt.Println(version)
	}
	if len(current) > 0 && err == nil {
		fmt.Println("current version: " + current)
	} else {
		fmt.Println("current version not available")
	}

	return nil
}

func retrieveVersions() ([]string, error) {
	versions := make([]string, 0)
	nvmcVersionsDir, err := util.GetVersionsPath()
	if err != nil {
		return versions, err
	}

	dirList, err := os.ReadDir(nvmcVersionsDir)
	if err != nil {
		return versions, err
	}

	for _, dirEntry := range dirList {
		if dirEntry.IsDir() && strings.HasPrefix(dirEntry.Name(), "v") {
			versions = append(versions, dirEntry.Name())
		}
	}

	return versions, nil
}

func currentVersion() (string, error) {
	nodeSymLink, err := util.GetSymLinkPath()
	if err != nil {
		return "", err
	}
	version, err := os.Readlink(nodeSymLink)
	// If the symlink or symlink destination does not exist then the current version is missing, but it isn't an error.
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	pathParts := strings.Split(version, string(os.PathSeparator))

	return pathParts[len(pathParts)-2], nil
}
