package file

import (
	"io"
	"os"
)

const MEGABYTE = 1024 * 1024
const DEFAULT_BUFFER_SIZE = 4 * MEGABYTE

type CopyOptions struct {
	ReplaceFile bool
	BufferSize  uint
}

func CopyFile(srcPath string, destPath string, opt *CopyOptions) error {
	if !opt.ReplaceFile && exists(destPath) {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
