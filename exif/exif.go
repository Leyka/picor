package exif

import (
	"fmt"
	"os"

	"github.com/Leyka/picor/geocoder"
	"github.com/rwcarlsen/goexif/exif"
)

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

func ExtractExif(filePath string, opt *ExifOptions) (*Exif, error) {
	if opt == nil {
		opt = DEFAULT_EXIF_OPTIONS
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decodedExif, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	var month, year string
	if opt.IncludeDate {
		month, year = extractMonthYear(decodedExif)
	}

	var location *geocoder.Location
	if opt.IncludeLocation {
		location = extractLocation(decodedExif)
	}

	return &Exif{
		Month:   month,
		Year:    year,
		City:    location.City,
		Country: location.Country,
	}, nil
}

func extractMonthYear(decodedExif *exif.Exif) (string, string) {
	dt, err := decodedExif.DateTime()
	if err != nil {
		return "", ""
	}

	month := dt.Format("01")
	year := dt.Format("2006")

	return month, year
}

func extractLocation(decodedExif *exif.Exif) *geocoder.Location {
	lat, long, err := decodedExif.LatLong()
	if err != nil {
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
