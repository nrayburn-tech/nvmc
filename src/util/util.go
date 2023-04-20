package util

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var VERSION = "UNSET"

func GetNvmcHomePath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	nvmcHome := filepath.Join(userHome, ".nvmc")

	return nvmcHome, nil
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

func GetVersionExePath(version string) (string, error) {
	versionDir, err := GetVersionPath(version)
	if err != nil {
		return "", err
	}
	var exeDir string
	if runtime.GOOS == "windows" {
		exeDir = versionDir
	} else {
		exeDir = filepath.Join(versionDir, "bin")
	}

	return exeDir, nil
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
