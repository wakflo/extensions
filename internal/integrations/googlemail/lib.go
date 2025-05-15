package googlemail

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googlemail/actions"
	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/extensions/internal/integrations/googlemail/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleMail())

type GoogleMail struct{}

func (n *GoogleMail) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *GoogleMail) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedGmailAuth,
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
