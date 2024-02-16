package cmd

import (
	"fmt"
	"nvmc/util"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("could not get current dir: %v", err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(dir, "nvmc")

	os.Exit(m.Run())
}

func runBinary(args ...string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	return cmd.CombinedOutput()
}

func TestInstall(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"install 18.2.0", []string{"install", "18.2.0"}, "there is not a current node version activated, will activate v18.2.0\nnow using node v18.2.0\nsuccessfully installed v18.2.0\n"},
		{"install 18.2.0 and use", []string{"install", "18.2.0", "--use"}, "there is not a current node version activated, will activate v18.2.0\nnow using node v18.2.0\nsuccessfully installed v18.2.0\n"},
	}

	for _, tt := range tests {
		util.IntegrationTest(t)
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args...)
			if err != nil {
				t.Fatalf("Output:%v\nError:%v", string(output), err)
			}

			actual := string(output)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("actual = %s, expected = %s", actual, tt.expected)
			}
		})
	}
}
