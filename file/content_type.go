package file

import (
	"net/http"
	"os"
	"regexp"
)

type IsTypeFn func(string) bool

func getFileContentType(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buffer [512]byte
	to, err := f.Read(buffer[:])
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:to])
	return contentType, nil
}

func IsTypeImageOrVideo(fileContentType string) bool {
	match, err := regexp.MatchString(`(?i)^(image|video)\/.*$`, fileContentType)
	if err != nil {
		return false
	}
	return match
}

func IsTypeImage(fileContentType string) bool {
	match, err := regexp.MatchString(`(?i)^image\/.*$`, fileContentType)
	if err != nil {
		return false
	}
	return match
}

func IsTypeVideo(fileContentType string) bool {
	match, err := regexp.MatchString(`(?i)^video\/.*$`, fileContentType)
	if err != nil {
		return false
	}
	return match
}
