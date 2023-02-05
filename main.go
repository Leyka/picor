package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/Leyka/picor/cache"
	"github.com/Leyka/picor/exif"
	"github.com/Leyka/picor/file"
	"github.com/Leyka/picor/geocoder"
	"github.com/joho/godotenv"
)

func setup() {
	// Loads .env in memory
	godotenv.Load()

	cache.SetupCache()

	geocoder.SetupGeocoderAPI()

	exif.Setup()
}

func Cleanup() {
	exif.CleanUp()
}

func main() {
	start := time.Now()
	setup()

	srcDir := "debug"
	destDir := "dest"

	filePaths, err := file.ListFilePathsByContentType(srcDir, file.IsTypeImageOrVideo)
	if err != nil {
		log.Panicln(err)
	}

	for _, srcPath := range filePaths {
		exif, err := exif.ExtractExif(srcPath, nil)
		if err != nil {
			log.Panicln(err)
		}

		destDirPath := filepath.Join(destDir, exif.Year, exif.Country, exif.City)
		file.CreateDirectoryIfNotExist(destDirPath)

		destDirFile := filepath.Join(destDirPath, filepath.Base(srcPath))
		err = file.CopyFile(srcPath, destDirFile, &file.CopyOptions{
			ReplaceFile: true,
			BufferSize:  file.DEFAULT_BUFFER_SIZE,
		})

		if err != nil {
			log.Panicln(err)
		}
	}

	Cleanup()
	fmt.Println("Done! Elapsed time: ", time.Since(start))
}
