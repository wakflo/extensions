package shared

import (
	"regexp"

	"google.golang.org/api/gmail/v1"

	"github.com/juicycleff/smartform/v1"
)

var (
	gmailForm = smartform.NewAuthForm("gmail-auth", "Gmail OAuth", smartform.AuthStrategyOAuth2)
	_         = gmailForm.
			OAuthField("oauth", "Gmail OAuth").
			AuthorizationURL("https://accounts.google.com/o/oauth2/auth").
			TokenURL("https://oauth2.googleapis.com/token").
			Scopes([]string{
			"https://mail.google.com/",
		}).
		Build()
)

var SharedGmailAuth = gmailForm.Build()

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
