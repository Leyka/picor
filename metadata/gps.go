package metadata

import (
	"fmt"

	"github.com/barasher/go-exiftool"
)

type GPSCoordinates struct {
	Latitude  float64
	Longitude float64
}

func extractGPSCoordinates(metadata *exiftool.FileMetadata) *GPSCoordinates {
	gpsLatitude := fmt.Sprintf("%v", metadata.Fields["GPSLatitude"])
	gpsLongitude := fmt.Sprintf("%v", metadata.Fields["GPSLongitude"])

	var lat, long float64
	fmt.Sscanf(gpsLatitude, "%f", &lat)
	fmt.Sscanf(gpsLongitude, "%f", &long)

	return &GPSCoordinates{
		Latitude:  lat,
		Longitude: long,
	}
}
