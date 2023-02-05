package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/Leyka/picor/cache"
	"github.com/Leyka/picor/exif"
	"github.com/Leyka/picor/file"
	"github.com/Leyka/picor/geocoder"
	"github.com/joho/godotenv"

	_ "net/http/pprof"
)

const MAX_WORKERS = 4

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
	defer Cleanup()

	srcDir := "debug"
	destDir := "dest"

	files, err := file.ListFilePathsByContentType(srcDir, file.IsTypeImageOrVideo)
	if err != nil {
		log.Panicln(err)
	}

	var wg sync.WaitGroup
	filesChan := make(chan string, MAX_WORKERS)

	for i := 0; i < MAX_WORKERS; i++ {
		wg.Add(1)
		go worker(filesChan, destDir, &wg)
	}

	for _, file := range files {
		filesChan <- file
	}

	close(filesChan)
	wg.Wait()
	fmt.Println("Done! Time elapsed:", time.Since(start))
}

func worker(srcFilePathsChan <-chan string, destDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Every time srcFilePathsChan is filled, this loop will be executed
	for srcPath := range srcFilePathsChan {
		exif, err := exif.ExtractExif(srcPath, nil)
		if err != nil {
			log.Println("Error extracting exif data from file:", err)
			continue
		}

		destDirPath := filepath.Join(destDir, exif.Year, exif.Country, exif.City)
		file.CreateDirectoryIfNotExist(destDirPath)

		destDirFile := filepath.Join(destDirPath, filepath.Base(srcPath))
		err = file.CopyFile(srcPath, destDirFile, &file.CopyOptions{
			ReplaceFile: true,
			BufferSize:  file.DEFAULT_BUFFER_SIZE,
		})

		if err != nil {
			log.Println("Error copying file:", err)
			continue
		}
	}
}
