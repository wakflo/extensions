package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var (
	form = smartform.NewAuthForm("convertKit-auth", "ConvertKit API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "ConvertKit API Key (Required)").
		Required(true).
		HelpText("The base URL of your app (e.g https://mystore.com) and it should start with HTTPS only")

	_ = form.TextField("api-secret", "API Secret (Required)").
		Required(true).
		HelpText("ConvertKit API Secret")

	ConvertKitSharedAuth = form.Build()
)
