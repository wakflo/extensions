package googlecalendar

import (
	"github.com/wakflo/go-sdk/autoform"
	"regexp"
	"time"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	sharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://www.googleapis.com/auth/calendar",
	}).Build()
)

// Function to clean the time string and return RFC3339 formatted time
func formatTimeString(input string) string {
	// Regular expression to strip the timezone identifier (e.g., [Africa/Lagos])
	re := regexp.MustCompile(`\[[A-Za-z/_]+\]`)
	cleanedString := re.ReplaceAllString(input, "")

	// Parse the cleaned time string (which should now be RFC3339 compatible)
	parsedTime, err := time.Parse(time.RFC3339, cleanedString)
	if err != nil {
		// Return empty string or any default value in case of error
		return "invalid date format"
	}

	// Return the time formatted back to RFC3339
	return parsedTime.Format(time.RFC3339)
}
