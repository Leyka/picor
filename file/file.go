package file

import (
	"io"
	"os"
	"path/filepath"
)

const MEGABYTE = 1024 * 1024
const DEFAULT_COPY_BUFFER_SIZE = 4 * MEGABYTE

func ListFiles(fromPath string, matchesType MatchesTypeFn) ([]string, error) {
	var files []string

	err := filepath.Walk(fromPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		contentType, err := getFileContentType(path)
		if err != nil {
			return err
		}
		if matchesType(contentType) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func CreateDirectoryIfNotExist(dirPath string) error {
	if !Exists(dirPath) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

type CopyOptions struct {
	ReplaceFile bool
	BufferSize  uint
}

func CopyFile(srcPath string, destPath string, opt *CopyOptions) error {
	if !opt.ReplaceFile && Exists(destPath) {
		return nil
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buffer := make([]byte, opt.BufferSize)
	for {
		nbReadBytes, err := srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if nbReadBytes == 0 {
			break
		}

		if _, err := destFile.Write(buffer[:nbReadBytes]); err != nil {
			return err
		}
	}

	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
