package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var (
	authForm = smartform.NewAuthForm("freshworkscrm-auth", "Freshworks CRM Auth", smartform.AuthStrategyCustom)
	_        = authForm.
			TextField("domain", "Freshworks Domain").
			Required(true).
			HelpText("The domain name of the freshworks account. eg. xyz.freshworks.com, type in only 'xyz'").
			Build()
)

var _ = authForm.
	TextField("api-key", "API Key").
	Required(true).
	HelpText("Your Freshworks CRM API key").
	Build()

var FreshworksCRMSharedAuth = authForm.Build()
