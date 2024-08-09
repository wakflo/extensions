package googlemail

import (
	"regexp"

	"google.golang.org/api/gmail/v1"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	sharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://mail.google.com/",
	}).Build()
)

var viewMailFormat = []*sdkcore.AutoFormSchema{
	{Const: "full", Title: "Full"},
	{Const: "minimal", Title: "Minimal"},
	{Const: "raw", Title: "Raw"},
	{Const: "metadata", Title: "Metadata"},
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func getHeader(headers []*gmail.MessagePartHeader, name string) string {
	for _, header := range headers {
		if header.Name == name {
			return header.Value
		}
	}
	return ""
}
