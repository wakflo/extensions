package googlemail

import (
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
