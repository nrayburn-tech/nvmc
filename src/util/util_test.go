package util

import (
	"testing"
)

func TestNormalizeVersionPrependsV(t *testing.T) {
	initVersion := "18.12.0"
	expectVersion := "v18.12.0"
	version, err := NormalizeVersion(initVersion)
	if err != nil || version != expectVersion {
		t.Fatalf(`NormalizeVersion(%q) = %q, %v, Wanted = %q`, initVersion, version, err, expectVersion)
	}
}

func TestNormalizeVersionLowercaseV(t *testing.T) {
	initVersion := "V18.12.0"
	expectVersion := "v18.12.0"
	version, err := NormalizeVersion(initVersion)
	if err != nil || version != expectVersion {
		t.Fatalf(`NormalizeVersion(%q) = %q, %v, Wanted = %q`, initVersion, version, err, expectVersion)
	}
}

func TestNormalizeVersionErrorOnEmptyVersion(t *testing.T) {
	initVersion := ""
	expectVersion := ""
	expectError := "version is required"
	version, err := NormalizeVersion(initVersion)
	if err.Error() != "version is required" || version != expectVersion {
		t.Fatalf(`NormalizeVersion(%q) = %q, %v, Wanted = %q, %v`, initVersion, version, err, expectVersion, expectError)
	}
}
