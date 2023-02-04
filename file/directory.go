package file

import (
	"os"
	"path/filepath"
)

// Recursively list all files in a directory by their content type (mime)
func ListFilePathsByContentType(rootPath string, isType IsTypeFn) ([]string, error) {
	var paths []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileContentType, err := getFileContentType(path)
		if err != nil {
			return err
		}

		if isType(fileContentType) {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

// Example: "path/to/dir" will create the 3 directories
func CreateDirectoryIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}
