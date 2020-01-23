package cmd


import (
	"testing"
)

func TestParseTimestamp(t *testing.T) {
	testCases := []string{
		"2020-01-01",
		"2020-12-31",
	}
	for _, tc := range testCases {
		_, err := parseTimestamp(tc)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

