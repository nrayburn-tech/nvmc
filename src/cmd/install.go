package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"nvmc/util"
	"os"
	"path/filepath"
	"strings"
)

type installCmd struct {
	command     *cobra.Command
	installOpts installOpts
	globalOpts  globalOpts
}

func newInstallCmd(globalOpts globalOpts) *installCmd {
	cmd := &installCmd{}
	cmd.command = &cobra.Command{
		Use:   "install <version>",
		Short: "Download and install <version>.",
		Example: `# Install version 18.2.0 and set it as active.
$ nvmc install 18.2.0 --use`,
		Args: cobra.ExactArgs(1),
		RunE: cmd.run(),
	}

	cmd.globalOpts = globalOpts
	cmd.command.Flags().BoolVar(&cmd.installOpts.overrideExistingInstall, "override-existing-install", defaultInstallOpts.overrideExistingInstall, "Overwrite existing installation.")
	cmd.command.Flags().BoolVar(&cmd.installOpts.skipChecksumValidation, "skip-checksum-validation", defaultInstallOpts.skipChecksumValidation, "Skip checksum validation after downloading.")
	cmd.command.Flags().BoolVar(&cmd.installOpts.use, "use", defaultInstallOpts.use, "After installing, set the installed <version> as active. (same as: nvmc use <version>).")

	return cmd
}

func (c *installCmd) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return install(version, c.globalOpts, c.installOpts)
	}
}

func install(version string, globalOpts globalOpts, installOpts installOpts) error {
	version, err := util.NormalizeVersion(version)
	if err != nil {
		return err
	}

	installationInfo, err := util.GetInstallationInfo(version)
	if err != nil {
		return err
	}

	// TODO: Validate the version
	versionDir, err := util.GetVersionPath(version)
	if err != nil {
		return err
	}

	if installOpts.overrideExistingInstall {
		if err := os.RemoveAll(versionDir); err != nil {
			return err
		}
	}

	if _, err := os.Stat(versionDir); err == nil {
		return errors.New("requested installation " + version + " already exists run with --override-existing-install to overwrite the existing version")
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	tempDir, err := os.MkdirTemp("", "nvmc-temp-"+version)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	tempZipFile, err := os.Create(filepath.Join(tempDir, installationInfo.FileNameWithExtension))
	if err != nil {
		return err
	}
	defer os.Remove(tempZipFile.Name())

	if err := util.Download(globalOpts.downloadUrl+"/"+version+"/"+installationInfo.FileNameWithExtension, tempZipFile); err != nil {
		return err
	}
	if _, err := tempZipFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if !installOpts.skipChecksumValidation {
		tempChecksumFile, err := os.Create(filepath.Join(tempDir, "SHASUMS256.txt"))
		if err != nil {
			return err
		}
		defer os.Remove(tempChecksumFile.Name())

		if err := util.Download(globalOpts.downloadUrl+"/"+version+"/SHASUMS256.txt", tempChecksumFile); err != nil {
			return err
		}
		if _, err := tempChecksumFile.Seek(0, io.SeekStart); err != nil {
			return err
		}

		fileBuf := new(bytes.Buffer)
		if _, err = fileBuf.ReadFrom(tempChecksumFile); err != nil {
			return err
		}
		fileContents := fileBuf.String()

		checksums := strings.Split(fileContents, "\n")
		for _, checksumLine := range checksums {
			if strings.HasSuffix(checksumLine, installationInfo.FileNameWithExtension) {
				checksum, found := strings.CutSuffix(checksumLine, " "+installationInfo.FileNameWithExtension)
				if !found {
					return errors.New("unable to verify checksum")
				}
				hash := sha256.New()
				if _, err := io.Copy(hash, tempZipFile); err != nil {
					return err
				}
				if _, err := tempZipFile.Seek(0, io.SeekStart); err != nil {
					return err
				}
				generatedChecksum := hex.EncodeToString(hash.Sum(nil))
				if strings.TrimSpace(checksum) != strings.TrimSpace(generatedChecksum) {
					return errors.New("checksum does not match")
				}
			}
		}
	}

	_, err = util.Unzip(tempZipFile, tempDir)
	if err != nil {
		return err
	}

	versionsDir, err := util.GetVersionsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(versionsDir, fs.ModePerm); err != nil {
		return err
	}
	if err := os.Rename(tempDir, versionDir); err != nil {
		return err
	}

	if _, err := currentVersion(); err != nil {
		fmt.Println("there is not a current node version activated, will activate " + version)
		installOpts.use = true
	}

	if installOpts.use {
		if err := use(version); err != nil {
			return err
		}
	}

	fmt.Printf("successfully installed %s\n", version)
	return nil
}
