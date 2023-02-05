package file

import (
	"os"
	"path/filepath"
)

func ListFiles(fromPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(fromPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
