package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var authForm = smartform.NewAuthForm("freshworkscrm-auth", "Freshworks CRM Auth", smartform.AuthStrategyCustom)
var _ = authForm.
	CustomField("domain", "Freshworks Domain").
	Required(true).
	HelpText("The domain name of the freshworks account. eg. xyz.freshworks.com, type in only 'xyz'").
	Build()

var _ = authForm.
	CustomField("api-key", "API Key").
	Required(true).
	HelpText("Your Freshworks CRM API key").
	Build()

var (
	FreshworksCRMSharedAuth = authForm.Build()
)
