package metadata

import (
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

	month, year := getMonthYear(&metadata)

	var city, country string = "", ""
	if options.FetchLocation {
		city, country = getLocation(&metadata)
	}

	return &Metadata{
		FilePath:     file,
		CreatedYear:  year,
		CreatedMonth: month,
		City:         city,
		Country:      country,
	}, nil
}

func getMonthYear(metadata *exiftool.FileMetadata) (string, string) {
	createDate := metadata.Fields["CreateDate"]
	if createDate == nil {
		return "", ""
	}

	// Format 2019:09:01 12:00:00
	month := createDate.(string)[5:7]
	year := createDate.(string)[0:4]
	return month, year
}

func getLocation(metadata *exiftool.FileMetadata) (string, string) {
	// TODO
	return "", ""
}
