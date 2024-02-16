package cmd

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"nvmc/util"
	"os"
	"sort"
	"strings"
)

type listCmd struct {
	command    *cobra.Command
	globalOpts globalOpts
	listOpts   listOpts
}

func newListCmd(globalOpts globalOpts) *listCmd {
	cmd := &listCmd{}
	cmd.command = &cobra.Command{
		Aliases: []string{"ls"},
		Use:     "list",
		Short:   "List all installed node versions.",
		Example: `$ nvmc list`,
		Args:    cobra.ExactArgs(0),
		RunE:    cmd.run(),
	}

	cmd.globalOpts = globalOpts

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

	current, _ := currentVersion()

	for _, version := range versions {
		if current == version {
			version = version + " (current)"
		}
		fmt.Println(version)
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

	semverVersions := make([]*semver.Version, 0)
	parseFailures := make([]string, 0)
	for _, version := range versions {
		semverVersion, err := semver.NewVersion(version)
		if err != nil {
			parseFailures = append(parseFailures, version)
		} else {
			semverVersions = append(semverVersions, semverVersion)
		}
	}

	sort.Sort(semver.Collection(semverVersions))

	versions = make([]string, len(semverVersions)+len(parseFailures))
	for i, semverVersion := range semverVersions {
		versions[i] = "v" + semverVersion.String()
	}
	for i, parseFailure := range parseFailures {
		versions[i+len(semverVersions)] = "Unable to parse " + parseFailure
	}

	return versions, nil
}

func currentVersion() (string, error) {
	nodeSymLink, err := util.GetSymLinkPath()
	if err != nil {
		return "", err
	}
	version, err := os.Readlink(nodeSymLink)
	if errors.Is(err, os.ErrNotExist) {
		return "", errors.New("Current version is missing. Version: " + version)
	} else if err != nil {
		return "", err
	}
	pathParts := strings.Split(version, string(os.PathSeparator))

	return pathParts[len(pathParts)-2], nil
}
