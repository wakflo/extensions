package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var (
	authForm = smartform.NewAuthForm("freshdesk-auth", "Freshdesk Auth", smartform.AuthStrategyCustom)
	_        = authForm.
			TextField("domain", "Freshdesk Domain").
			Required(true).
			HelpText("The domain name of the freshdesk account. eg. xyz.freshdesk.com, type in only 'xyz'").
			Build()
)

var _ = authForm.
	TextField("api-key", "API Key").
	Required(true).
	HelpText("The api key used to authenticate freshdesk.").
	Build()

var FreshdeskSharedAuth = authForm.Build()
