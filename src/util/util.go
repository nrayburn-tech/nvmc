package util

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// VERSION nvmc library version set at compile time in release.yaml
var VERSION = "UNSET"

type InstallationInfo struct {
	FileExtension            string
	FileNameWithExtension    string
	FileNameWithoutExtension string
}

func GetInstallationInfo(version string) (*InstallationInfo, error) {
	version, err := NormalizeVersion(version)
	if err != nil {
		return nil, err
	}

	normalizedArch := getNodeArch()
	normalizedPlatform := getNodeOs()

	fileExtension := getFileExtension()
	fileNameWithoutExtension := "node-" + version + "-" + normalizedPlatform + "-" + normalizedArch
	fileNameWithExtension := fileNameWithoutExtension + fileExtension

	return &InstallationInfo{fileExtension, fileNameWithExtension, fileNameWithoutExtension}, nil
}

func GetNvmcHomePath() (string, error) {
	nvmcHome := os.Getenv("NVMC_HOME")
	var home string
	if len(nvmcHome) == 0 {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		home = filepath.Join(userHome, ".nvmc")
	} else {
		home = filepath.Clean(nvmcHome)
	}

	return home, nil
}

func GetVersionsPath() (string, error) {
	nvmcHome, err := GetNvmcHomePath()
	if err != nil {
		return "", err
	}
	nvmcVersionsDir := filepath.Join(nvmcHome, "versions")

	return nvmcVersionsDir, nil
}

func GetVersionPath(version string) (string, error) {
	versionsDir, err := GetVersionsPath()
	if err != nil {
		return "", err
	}
	versionDir := filepath.Join(versionsDir, version)

	return versionDir, nil
}

func GetSymLinkPath() (string, error) {
	nvmcHome, err := GetNvmcHomePath()
	if err != nil {
		return "", err
	}
	symLink := filepath.Join(nvmcHome, "nodejs")
	return symLink, nil
}

func NormalizeVersion(version string) (string, error) {
	if len(version) == 0 {
		return "", errors.New("version is required")
	}
	version = strings.ToLower(version)
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	return version, nil
}

func getFileExtension() string {
	if runtime.GOOS == "windows" {
		return ".zip"
	}

	return ".tar.gz"
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
