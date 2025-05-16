// shared/auth.go
package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var (
	form = smartform.NewAuthForm("sendowl-auth", "Send Owl API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api_key", "API Key").
		Required(true).
		HelpText("Your SendOwl API key found in your account settings")

	_ = form.TextField("api_secret", "API Secret").
		Required(true).
		HelpText("Your SendOwl API secret found in your account settings")

	SharedAuth = form.Build()
)
