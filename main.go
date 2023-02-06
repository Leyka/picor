package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Leyka/picor/cache"
	"github.com/Leyka/picor/file"
	"github.com/Leyka/picor/metadata"
	"github.com/schollz/progressbar/v3"
)

var NUM_WORKERS = runtime.NumCPU()

func setup() {
	cache.Setup()
	metadata.Setup(NUM_WORKERS)
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
	fmt.Println("~ Found", totalFiles, "images/videos to organize")

	// Create a channel to keep track of the number of processed files
	processedFilesChan := make(chan int, NUM_WORKERS)

	go startWorkers(files, processedFilesChan, "dest")

	bar := progressbar.NewOptions(totalFiles,
		progressbar.OptionSetDescription("copying"),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
			fmt.Println("~ Completed in", time.Since(start))
		}),
	)

	processedFiles := 0
	for count := range processedFilesChan {
		// Update progress bar
		bar.Add(count)

		processedFiles += count
		if processedFiles == totalFiles {
			close(processedFilesChan)
			bar.Finish()
		}
	}
}

func startWorkers(files []string, processedFilesChan chan<- int, destDir string) {
	filesChan := make(chan string, NUM_WORKERS)
	defer close(filesChan)

	// Start workers that will receive files to process
	for i := 0; i < NUM_WORKERS; i++ {
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

		// TODO: Handle when no date, no country etc...
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
