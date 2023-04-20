package util

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(zipFile fs.File, basePath string) (string, error) {
	fileInfo, err := zipFile.Stat()
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(fileInfo.Name(), ".zip") {
		return ZipUnzip(zipFile, basePath)
	} else {
		return TarGzUnzip(zipFile, basePath)
	}
}

func TarGzUnzip(zipFile fs.File, basePath string) (string, error) {
	zipReadCloser, err := gzip.NewReader(zipFile)
	if err != nil {
		return "", err
	}
	defer zipReadCloser.Close()

	tarReadCloser := tar.NewReader(zipReadCloser)

	header, err := tarReadCloser.Next()
	unzippedFilePath := filepath.Join(basePath, header.Name)

	for true {
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return "", err
		}
		path := filepath.Join(basePath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return "", err
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return "", err
			}

			if err := os.Symlink(header.Linkname, path); err != nil {
				return "", err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return "", err
			}
			outFile, err := os.Create(path)
			if err != nil {
				return "", err
			}

			if err = os.Chmod(path, fs.FileMode(header.Mode)); err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tarReadCloser); err != nil {
				return "", err
			}
			// Don't defer close to avoid buffering everything in memory.
			outFile.Close()
		default:
			return "", errors.New(fmt.Sprintf("unsupported tar type: %v for name %s", header.Typeflag, header.Name))
		}

		header, err = tarReadCloser.Next()
	}

	return unzippedFilePath, nil
}

func ZipUnzip(zipFile fs.File, basePath string) (string, error) {
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, zipFile)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(buff.Bytes())
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return "", err
	}

	unzippedFilePath := filepath.Join(basePath, zipReader.File[0].Name)

	for _, zipFile := range zipReader.File {
		// TODO: Close fileReader
		fileReader, err := zipFile.Open()
		if err != nil {
			return "", err
		}
		path := filepath.Join(basePath, zipFile.Name)

		if zipFile.FileInfo().IsDir() {

			if err := os.MkdirAll(path, zipFile.Mode()); err != nil {
				return "", err
			}
		} else {

			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return "", err
			}
			// TODO: Close file
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
			if err != nil {
				return "", err
			}

			if _, err = io.Copy(file, fileReader); err != nil {
				return "", err
			}
		}
	}
	return unzippedFilePath, nil
}
