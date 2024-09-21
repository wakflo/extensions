package calendly

import (
	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL   = "https://auth.calendly.com/oauth/token"
	sharedAuth = autoform.NewOAuthField("https://api.calendly.com", &tokenURL, []string{}).Build()
)
