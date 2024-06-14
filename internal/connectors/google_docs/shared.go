package googledocs

import (
	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	sharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://www.googleapis.com/auth/documents",
	}).Build()
)
