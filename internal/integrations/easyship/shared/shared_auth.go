package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var (
	authForm = smartform.NewAuthForm("easyship-auth", "EasyShip Auth", smartform.AuthStrategyCustom)
	_        = authForm.
			CustomField("api-key", "API Key").
			Required(true).
			HelpText("API Application Key").
			Build()
)

var EasyShipSharedAuthV2 = authForm.Build()
