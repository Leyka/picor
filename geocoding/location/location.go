package location

type Location struct {
	Country     string
	CountryCode string
	State       string
	City        string
}

var UNKNOWN_LOCATION = &Location{
	Country:     "~ unknown country",
	CountryCode: "~ unknown country code",
	State:       "~ unknown state",
	City:        "~ unknown city",
}
