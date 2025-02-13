package notion

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/notion/actions"
	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/extensions/internal/integrations/notion/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewNotion(), Flow, ReadME)

type Notion struct{}

func (n *Notion) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
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
