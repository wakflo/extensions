package notion

import (
	"github.com/wakflo/extensions/internal/integrations/notion/actions"
	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/extensions/internal/integrations/notion/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewNotion())

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
