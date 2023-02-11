package metadata

import (
	"fmt"

	"github.com/Leyka/picor/geocoding"
	"github.com/barasher/go-exiftool"
)

type Options struct {
	FetchLocation bool
}

type Metadata struct {
	FilePath     string
	CreatedYear  string
	CreatedMonth string
	City         string
	Country      string
}

var instances []*exiftool.Exiftool

func Setup(nbInstances int) {
	for i := 0; i < nbInstances; i++ {
		instance, err := exiftool.NewExiftool(exiftool.CoordFormant("%+.6f")) // signed coordinates
		if err != nil {
			panic(err)
		}
		instances = append(instances, instance)
	}
}

func Close() {
	for _, instance := range instances {
		instance.Close()
	}
}

func ExtractMetadata(instanceId int, file string, options *Options) (*Metadata, error) {
	instance := instances[instanceId]
	if instance == nil {
		panic("Instance not found You must run metadata.Setup()")
	}

	metadata := instance.ExtractMetadata(file)[0]
	if metadata.Err != nil {
		return &Metadata{
			FilePath: file,
		}, metadata.Err
	}

	date := extractDate(&metadata)
	if date == nil {
		// Fallback to file name
		date = tryGetDateFromFile(file)
		if date == nil {
			date = &Date{
				year:  "",
				month: "",
				day:   "",
			}
		}
	}

	var city, country string = "", ""
	if options.FetchLocation {
		coord := extractGPSCoordinates(&metadata)
		if coord.Latitude > 0 && coord.Longitude > 0 {
			location, err := geocoding.ReverseGeocoding(coord.Latitude, coord.Longitude)
			if err == nil {
				city = location.City
				country = location.Country
			} else {
				// TODO: Silent Log error
				fmt.Println(err, coord)
			}
		}
	}

	return &Metadata{
		FilePath:     file,
		CreatedYear:  date.year,
		CreatedMonth: date.month,
		City:         city,
		Country:      country,
	}, nil
}
