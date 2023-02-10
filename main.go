package main

import (
	"os"
	"runtime"

	"github.com/Leyka/picor/geocoding"
	"github.com/joho/godotenv"
)

var NUM_WORKERS = runtime.NumCPU()

func main() {
	setup()
	defer cleanup()

}

func setup() {
	godotenv.Load()

	SetupExiftool(NUM_WORKERS)

	geocoding.Setup(geocoding.GeocodingSettings{
		TomTomApiKey: os.Getenv("TOMTOM_API_KEY"),
	})
}

func cleanup() {
	CloseExiftool()

	geocoding.Close()
}
