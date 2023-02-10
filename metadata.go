package main

import "github.com/barasher/go-exiftool"

type Metadata struct {
	FilePath     string
	CreatedYear  string
	CreatedMonth string
	City         string
	Country      string
}

// Pool of exiftool instances
var instances []*exiftool.Exiftool

func SetupExiftool(nbInstances int) {
	for i := 0; i < nbInstances; i++ {
		instance, err := exiftool.NewExiftool(exiftool.CoordFormant("%+.6f")) // signed coordinates
		if err != nil {
			panic(err)
		}
		instances = append(instances, instance)
	}
}

func CloseExiftool() {
	for _, instance := range instances {
		instance.Close()
	}
}
