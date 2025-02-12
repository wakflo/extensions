package monday

import (
	"github.com/wakflo/extensions/internal/integrations/monday/actions"
	"github.com/wakflo/extensions/internal/integrations/monday/shared"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewMonday())

type Monday struct{}

func (n *Monday) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Monday) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Monday) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateUpdateAction(),

		actions.NewCreateItemAction(),

		actions.NewCreateGroupAction(),

		actions.NewCreateColumnAction(),
	}
}

func NewMonday() sdk.Integration {
	return &Monday{}
}
