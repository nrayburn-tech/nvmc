package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"nvmc/util"
	"os"
	"runtime"
)

type installOptions struct {
	use bool
}

type installCmd struct {
	command *cobra.Command
	options installOptions
}

func newInstallCmd() *installCmd {
	cmd := &installCmd{}
	cmd.command = &cobra.Command{
		Use:   "install <version>",
		Short: "Download and install <version>.",
		Example: `# Install version 18.2.0 and set it as active.
$ nvmc install 18.2.0 --use`,
		Args: cobra.ExactArgs(1),
		RunE: cmd.run(),
	}

	cmd.command.Flags().BoolVar(&cmd.options.use, "use", true, "After installing, set the installed <version> as active. (same as: nvmc use <version>).")

	return cmd
}

func (c *installCmd) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return install(version, c.options)
	}
}

func install(version string, opts installOptions) error {
	version, err := util.NormalizeVersion(version)
	if err != nil {
		return err
	}

	// TODO: Validate the version

	versionDir, err := util.GetVersionPath(version)
	if err != nil {
		return err
	}
	tempDir, err := os.MkdirTemp("", "nvmc-"+version)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if _, err := os.Stat(versionDir); err == nil {
		return errors.New("requested installation " + version + " already exists, uninstall it and try again")
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	fileName := "node-" + version + "-" + getNodeOs() + "-" + getNodeArch() + getDownloadExtension()
	tempFile, err := os.CreateTemp(tempDir, "nvmc-"+fileName)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// TODO: Download and validate the checksum.
	if err := util.Download("https://nodejs.org/dist/"+version+"/"+fileName, tempFile); err != nil {
		return err
	}
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	destPath, err := util.Unzip(tempFile, tempDir)
	if err != nil {
		return err
	}

	if err := os.Rename(destPath, versionDir); err != nil {
		return err
	}

	if currVersion, err := currentVersion(); len(currVersion) == 0 && err == nil {
		fmt.Println("there is not a current node version activated, will activate " + version)
		opts.use = true
	}

	if opts.use {
		if err := use(version); err != nil {
			return err
		}
	}

	return nil
}

func getNodeOs() string {
	switch runtime.GOOS {
	case "windows":
		return "win"
	default:
		return runtime.GOOS
	}
}

func getNodeArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x64"
	case "386":
		return "x86"
	default:
		return runtime.GOARCH
	}
}

func getDownloadExtension() string {
	if runtime.GOOS == "windows" {
		return ".zip"
	}

	return ".tar.gz"
}
