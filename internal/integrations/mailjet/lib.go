package mailjet

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/mailjet/actions"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/extensions/internal/integrations/mailjet/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewMailJet())

type MailJet struct{}

func (n *MailJet) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (m *MailJet) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
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
