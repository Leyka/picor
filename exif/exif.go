package exif

import (
	"fmt"
	"log"

	"github.com/Leyka/picor/geocoder"
	"github.com/barasher/go-exiftool"
)

var instance *exiftool.Exiftool

type ExifOptions struct {
	IncludeDate     bool
	IncludeLocation bool
}

var DEFAULT_EXIF_OPTIONS = &ExifOptions{
	IncludeDate:     true,
	IncludeLocation: true,
}

type Exif struct {
	Month   string
	Year    string
	City    string
	Country string
}

func Setup() {
	// Create new instance of exiftool with 4 decimal places for lat/long
	et, err := exiftool.NewExiftool(exiftool.CoordFormant("%+.6f"))
	if err != nil {
		log.Panicln(err)
	}
	instance = et
}

func CleanUp() {
	instance.Close()
}

func ExtractExif(filePath string, opt *ExifOptions) (*Exif, error) {
	if opt == nil {
		opt = DEFAULT_EXIF_OPTIONS
	}

	metadata := instance.ExtractMetadata(filePath)[0]

	var month, year string
	if opt.IncludeDate {
		month, year = extractMonthYear(&metadata)
	}

	var location *geocoder.Location
	if opt.IncludeLocation {
		location = extractLocation(&metadata)
	}

	return &Exif{
		Month:   month,
		Year:    year,
		City:    location.City,
		Country: location.Country,
	}, nil
}

func extractMonthYear(metadata *exiftool.FileMetadata) (string, string) {
	createDate := metadata.Fields["CreateDate"]
	if createDate == nil {
		return "", ""
	}

	// Format 2019:09:01 12:00:00
	month := createDate.(string)[5:7]
	year := createDate.(string)[0:4]
	return month, year
}

func extractLocation(metadata *exiftool.FileMetadata) *geocoder.Location {
	gpsLatitude := fmt.Sprintf("%v", metadata.Fields["GPSLatitude"])
	gpsLongitude := fmt.Sprintf("%v", metadata.Fields["GPSLongitude"])

	var lat, long float64
	fmt.Sscanf(gpsLatitude, "%f", &lat)
	fmt.Sscanf(gpsLongitude, "%f", &long)

	if lat == 0 && long == 0 {
		return geocoder.UNKNOWN_LOCATION
	}

	geocode := geocoder.NewGeoCode(lat, long)
	location, err := geocode.GetLocation()
	if err != nil {
		fmt.Errorf("Error getting location from API: %s", err)
		return geocoder.UNKNOWN_LOCATION
	}

	return location
}
