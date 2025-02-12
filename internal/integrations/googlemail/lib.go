package googlemail

import (
	"github.com/wakflo/extensions/internal/integrations/googlemail/actions"
	"github.com/wakflo/extensions/internal/integrations/googlemail/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewGoogleMail())

type GoogleMail struct{}

func (n *GoogleMail) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *GoogleMail) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewEmailTrigger(),
	}
}

func (n *GoogleMail) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewSendEmailTemplateAction(),

		actions.NewSendEmailAction(),

		actions.NewListMailsAction(),

		actions.NewGetThreadAction(),

		actions.NewGetMailAction(),
	}
}

func NewGoogleMail() sdk.Integration {
	return &GoogleMail{}
}
