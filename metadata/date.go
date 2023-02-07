package metadata

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/barasher/go-exiftool"
)

type Date struct {
	year  string
	month string
	day   string
}

func extractDate(metadata *exiftool.FileMetadata) *Date {
	createDate := metadata.Fields["CreateDate"]
	if createDate == nil {
		return nil
	}

	// Format 2019:09:01 12:00:00
	return &Date{
		year:  createDate.(string)[0:4],
		month: createDate.(string)[5:7],
		day:   createDate.(string)[8:10],
	}
}

// Try to get the date from the file name
func tryGetDateFromFile(fileName string) *Date {
	formats := []string{
		`(\d{4})[\D]*(\d{2})[\D]*(\d{2})`, // YYYYMMDD (or YYYYDDMM...)
		`(\d{2})[\D]*(\d{2})[\D]*(\d{4})`, // DDMMYYYY (or MMDDYYYY)
	}

	for _, format := range formats {
		r := regexp.MustCompile(format)
		matches := r.FindStringSubmatch(fileName)

		if len(matches) == 4 {
			// Check which format we have
			if isMaybeYear(matches[1]) {
				// YYYYMMDD
				month := matches[2]
				day := matches[3]
				if !isMaybeMonth(month) {
					// YYYYDDMM
					month = matches[3]
					day = matches[2]
				}

				return &Date{
					year:  matches[1],
					month: month,
					day:   day,
				}
			} else if isMaybeYear(matches[3]) {
				// DDMMYYYY
				month := matches[2]
				day := matches[1]
				if !isMaybeMonth(month) {
					// MMDDYYYY
					month = matches[1]
					day = matches[2]
				}
				return &Date{
					year:  matches[3],
					month: month,
					day:   day,
				}
			}
		}
	}

	return nil
}

func isMaybeYear(maybeYear string) bool {
	if len(maybeYear) != 4 {
		return false
	}

	return strings.HasPrefix(maybeYear, "20") || strings.HasPrefix(maybeYear, "19")
}

func isMaybeMonth(maybeMonth string) bool {
	if len(maybeMonth) != 2 {
		return false
	}

	maybeMonthInt, err := strconv.Atoi(maybeMonth)
	if err != nil {
		return false
	}

	if maybeMonthInt > 12 {
		return false
	}

	// Still not sure if it's a day or a month at this point... It's a guess.
	return true
}
