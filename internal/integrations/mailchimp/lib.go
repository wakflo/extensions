package mailchimp

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/actions"
	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/extensions/internal/integrations/mailchimp/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewMailchimp())

type Mailchimp struct{}

func (n *Mailchimp) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Mailchimp) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (n *Mailchimp) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewUnsubscriberTrigger(),

		triggers.NewNewSubscriberTrigger(),
	}
}

func (n *Mailchimp) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateSubscriberStatusAction(),

		actions.NewRemoveSubscriberFromTagAction(),

		actions.NewGetAllListAction(),

		actions.NewAddSubscriberToTagAction(),

		actions.NewAddNoteToSubscriberAction(),

		actions.NewAddMemberToListAction(),
	}
}

func NewMailchimp() sdk.Integration {
	return &Mailchimp{}
}
