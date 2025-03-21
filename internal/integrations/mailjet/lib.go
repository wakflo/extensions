package mailjet

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/mailjet/actions"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/extensions/internal/integrations/mailjet/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewMailJet(), Flow, ReadME)

type MailJet struct{}

func (m *MailJet) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (m *MailJet) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewEmailSentTrigger(),
		triggers.NewContactUpdatedTrigger(),
	}
}

func (m *MailJet) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewSendEmailAction(),
		actions.NewGetContactAction(),
		actions.NewCreateContactAction(),
		actions.NewListContactsAction(),
	}
}

func NewMailJet() sdk.Integration {
	return &MailJet{}
}
