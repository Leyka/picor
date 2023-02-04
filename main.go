package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Leyka/picor/cache"
	"github.com/Leyka/picor/geocoder"
	"github.com/joho/godotenv"
	"github.com/rwcarlsen/goexif/exif"
)

type Exif struct {
	year     string
	location *geocoder.Location
}

func main() {
	godotenv.Load()
	cache.SetupCache()
	geocoder.SetupGeocoderAPI()

	srcDir := "debug"
	destDir := "dest"

	mediaPaths, err := listMediaPaths(srcDir)
	if err != nil {
		log.Panic(err)
	}

	for _, src := range mediaPaths {
		exif, err := extractExif(src)
		if err != nil {
			log.Fatal(err)
		}

		destYearDir := filepath.Join(destDir, exif.year, exif.location.Country, exif.location.City)
		createDestDir(destYearDir)

		dest := filepath.Join(destYearDir, filepath.Base(src))
		copyFile(src, dest, true)
	}

	log.Println("Done!")
}

func listMediaPaths(rootPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isImageType(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func isImageType(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	var buffer [512]byte
	to, err := f.Read(buffer[:])
	if err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer[:to])

	return contentType == "image/jpeg" ||
		contentType == "image/png" ||
		contentType == "image/gif" ||
		contentType == "image/heif" ||
		contentType == "image/heic"
}

func extractExif(path string) (*Exif, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decodedExif, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	dt, err := decodedExif.DateTime()
	if err != nil {
		return nil, err
	}
	year := dt.Format("2006")

	lat, long, err := decodedExif.LatLong()
	if err != nil {
		return &Exif{
			year: year,
			location: &geocoder.Location{
				Country: "Unknown Country",
				City:    "Unknown City",
			},
		}, nil
	}

	geocode := geocoder.NewGeoCode(lat, long)
	location, err := geocode.GetLocation()
	if err != nil {
		return nil, err
	}

	return &Exif{
		year:     year,
		location: location,
	}, nil
}

func createDestDir(destDir string) {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}
}

func copyFile(srcPath string, destPath string, replaceFile bool) {
	if !replaceFile && fileExists(destPath) {
		return
	}

	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("Failed to open source file: %v", err)
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		log.Fatalf("Failed to create destination file: %v", err)
	}
	defer dest.Close()

	// 5 MB chunks
	buf := make([]byte, 5*1024*1024)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("Failed to read source file: %v", err)
		}

		if n == 0 {
			break
		}

		if _, err := dest.Write(buf[:n]); err != nil {
			log.Fatalf("Failed to write destination file: %v", err)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
