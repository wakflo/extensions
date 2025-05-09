// shared/auth.go
package shared

import (
	"github.com/wakflo/go-sdk/autoform"
	"github.com/wakflo/go-sdk/core"
)

var SharedAuth = autoform.NewCustomAuthField().
	SetLabel("SendOwl API Authentication").
	SetDescription("Enter your SendOwl API key and secret to connect").
	SetRequired(true).
	SetFields(map[string]*core.AutoFormSchema{
		"api_key": autoform.NewShortTextField().
			SetLabel("API Key").
			SetDescription("Your SendOwl API key found in your account settings").
			SetRequired(true).
			Build(),
		"api_secret": autoform.NewShortTextField().
			SetLabel("API Secret").
			SetDescription("Your SendOwl API secret found in your account settings").
			SetRequired(true).
			Build(),
	}).
	Build()
