package zendeskapp

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/zendeskapp/actions"
	"github.com/wakflo/extensions/internal/integrations/zendeskapp/shared"
	"github.com/wakflo/extensions/internal/integrations/zendeskapp/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewZendeskApp())

type Zendesk struct{}

func (z *Zendesk) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (z *Zendesk) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (z *Zendesk) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewTicketCreatedTrigger(),
	}
}

func (z *Zendesk) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetGroupsAction(),
		actions.NewGetTicketsAction(),
	}
}

func NewZendeskApp() sdk.Integration {
	return &Zendesk{}
}
