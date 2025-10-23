// shared/auth.go
package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var (
	form = smartform.NewAuthForm("socialkit-auth", "SocialKit API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("access_key", "Access Key").
		Required(true).
		HelpText("Your SocialKit API access key found in your account settings")

	SharedAuth = form.Build()
)
