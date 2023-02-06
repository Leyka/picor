package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/Leyka/picor/cache"
	"github.com/Leyka/picor/file"
	"github.com/Leyka/picor/metadata"
	"github.com/schollz/progressbar/v3"
)

const MAX_WORKERS = 1

func setup() {
	cache.Setup()
	metadata.Setup(MAX_WORKERS)
}

func cleanup() {
	cache.Close()
	metadata.Close()
}

func main() {
	start := time.Now()

	setup()
	defer cleanup()

	files, err := file.ListFiles("test_photos", file.IsTypeImageOrVideo)
	if err != nil {
		log.Panicln(err)
	}
	totalFiles := len(files)
	fmt.Println("Found", totalFiles, "files to process")

	// Create a channel to keep track of the number of processed files
	processedFilesChan := make(chan int, MAX_WORKERS)

	go startWorkers(files, processedFilesChan, "dest")

	bar := progressbar.NewOptions(totalFiles,
		progressbar.OptionSetDescription("copying"),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	)

	processedFiles := 0
	for count := range processedFilesChan {
		// Update progress bar
		processedFiles += count
		bar.Add(processedFiles)

		if processedFiles == totalFiles {
			close(processedFilesChan)
		}
	}

	// All files have been processed
	bar.Finish()
	fmt.Println()
	fmt.Println("Done! It took", time.Since(start))
}

func startWorkers(files []string, processedFilesChan chan<- int, destDir string) {
	filesChan := make(chan string, MAX_WORKERS)
	defer close(filesChan)

	// Start workers that will receive files to process
	for i := 0; i < MAX_WORKERS; i++ {
		go processFileWorker(i, filesChan, processedFilesChan, destDir)
	}

	for _, file := range files {
		filesChan <- file
	}
}

func processFileWorker(id int, filesChan <-chan string, processedFilesChan chan<- int, destDir string) {
	for srcFile := range filesChan {
		metadata, err := metadata.ExtractMetadata(id, srcFile, &metadata.Options{
			FetchLocation: false,
		})
		if err != nil {
			// TODO: Silent log in file
			continue
		}

		destDirPath := filepath.Join(destDir, metadata.CreatedYear)
		file.CreateDirectoryIfNotExist(destDirPath)

		destDirFile := filepath.Join(destDirPath, filepath.Base(srcFile))

		err = file.CopyFile(srcFile, destDirFile, &file.CopyOptions{
			ReplaceFile: true,
			BufferSize:  file.DEFAULT_COPY_BUFFER_SIZE,
		})
		if err != nil {
			// TODO: Silent log in file
			continue
		}

		processedFilesChan <- 1
	}
}
