package trello

import (
	"github.com/wakflo/extensions/internal/integrations/trello/actions"
	"github.com/wakflo/extensions/internal/integrations/trello/shared"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewTrello())

type Trello struct{}

func (n *Trello) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Trello) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Trello) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateCardAction(),

		actions.NewFindCardAction(),

		actions.NewDeleteCardAction(),

		actions.NewCreateListAction(),

		actions.NewDeleteCardAction(),
	}
}

func NewTrello() sdk.Integration {
	return &Trello{}
}
