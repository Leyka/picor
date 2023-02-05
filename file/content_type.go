package file

import (
	"bytes"
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
