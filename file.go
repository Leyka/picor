package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const MEGABYTE = 1024 * 1024
const DEFAULT_COPY_BUFFER_SIZE = 4 * MEGABYTE

type MatchesTypeFn func(string) bool

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

func getFileContentType(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var header [512]byte
	to, err := f.Read(header[:])
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(header[:to])

	// Fix for heic files: http.DetectContentType returns "application/octet-stream"
	// Read the magic number of heic file
	if contentType == "application/octet-stream" && bytes.Equal(header[4:12], []byte("ftypheic")) {
		contentType = "image/heic"
	}

	return contentType, nil
}

func IsTypeImageOrVideo(fileContentType string) bool {
	return IsTypeImage(fileContentType) || IsTypeVideo(fileContentType)
}

func IsTypeImage(fileContentType string) bool {
	if !strings.HasPrefix(fileContentType, "image/") {
		return false
	}

	// TODO: Validate supported format by exiftool
	return fileContentType == "image/jpeg" ||
		fileContentType == "image/png" ||
		fileContentType == "image/heic" ||
		fileContentType == "image/heif" ||
		fileContentType == "image/tiff" ||
		fileContentType == "image/bmp" ||
		fileContentType == "image/gif"
}

func IsTypeVideo(fileContentType string) bool {
	if !strings.HasPrefix(fileContentType, "video/") {
		return false
	}

	// TODO: Validate supported format by exiftool
	return fileContentType == "video/mp4" ||
		fileContentType == "video/quicktime" ||
		fileContentType == "video/x-msvideo" ||
		fileContentType == "video/x-ms-wmv" ||
		fileContentType == "video/x-flv" ||
		fileContentType == "video/x-matroska"
}

func CreateDirectoryIfNotExist(dirPath string) error {
	if !exists(dirPath) {
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
