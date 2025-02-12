package shared

import (
	"regexp"

	"google.golang.org/api/gmail/v1"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	SharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://mail.google.com/",
	}).Build()
)

var ViewMailFormat = []*sdkcore.AutoFormSchema{
	{Const: "full", Title: "Full"},
	{Const: "minimal", Title: "Minimal"},
	{Const: "raw", Title: "Raw"},
	{Const: "metadata", Title: "Metadata"},
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func GetHeader(headers []*gmail.MessagePartHeader, name string) string {
	for _, header := range headers {
		if header.Name == name {
			return header.Value
		}
	}
	return ""
}
