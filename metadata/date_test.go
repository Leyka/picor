package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryGetDateFromFile(t *testing.T) {
	inputs := []string{
		"blablabla-20190901_201909-blabla.jpg",
		"2019-09-01_201909.jpg",
		"01-09_2019",
	}

	for _, input := range inputs {
		date := TryGetDateFromFile(input)
		assert.NotNil(t, date)
		assert.Equal(t, "2019", date.year)
		assert.Equal(t, "09", date.month)
		assert.Equal(t, "01", date.day)
	}

	// Special case
	input := "20120728_121137.jpg"
	date := TryGetDateFromFile(input)
	assert.NotNil(t, date)
	assert.Equal(t, "2012", date.year)
	assert.Equal(t, "07", date.month)
	assert.Equal(t, "28", date.day)
}
