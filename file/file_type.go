package file

import (
	"bytes"
	"net/http"
	"os"
	"strings"
)

type MatchesTypeFn func(string) bool

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
