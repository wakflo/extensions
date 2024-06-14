package googlesheets

import (
	"strconv"

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	sharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://www.googleapis.com/auth/spreadsheets",
	}).Build()
)

// this functon converts string to int64
func convertToInt64(s string) int64 {
	convertedString, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return convertedString
}
