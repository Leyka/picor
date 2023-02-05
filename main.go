package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Leyka/picor/cache"
	"github.com/Leyka/picor/file"
	"github.com/schollz/progressbar/v3"
)

const MAX_WORKERS = 1

func setup() {
	cache.Setup()
}

func cleanup() {
	cache.Close()
}

func main() {
	start := time.Now()

	setup()
	defer cleanup()

	files, err := file.ListFiles("test_photos")
	if err != nil {
		log.Panicln(err)
	}
	totalFiles := len(files)
	fmt.Println("Found", totalFiles, "files to process")

	// Create a channel to keep track of the number of processed files
	processedFilesChan := make(chan int, MAX_WORKERS)

	go startWorkers(files, processedFilesChan)

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

func startWorkers(files []string, processedFilesChan chan<- int) {
	filesChan := make(chan string, MAX_WORKERS)
	defer close(filesChan)

	// Start workers that will receive files to process
	for i := 0; i < MAX_WORKERS; i++ {
		go processFileWorker(i, filesChan, processedFilesChan)
	}

	for _, file := range files {
		filesChan <- file
	}
}

func processFileWorker(id int, filesChan <-chan string, processedFilesChan chan<- int) {
	for srcFile := range filesChan {
		// Do something with the file
		if srcFile == "" {
			// TODO
		}

		time.Sleep(1000 * time.Millisecond)

		processedFilesChan <- 1
	}
}
