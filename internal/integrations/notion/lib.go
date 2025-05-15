package notion

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/notion/actions"
	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/extensions/internal/integrations/notion/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewNotion())

type Notion struct{}

func (n *Notion) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Notion) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedNotionAuth,
	}
}

func (n *Notion) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewPageCreatedTrigger(),
	}
}

func (n *Notion) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdatePageAction(),

		actions.NewRetrievePageAction(),

		actions.NewCreatePageAction(),
	}
}

func NewNotion() sdk.Integration {
	return &Notion{}
}
