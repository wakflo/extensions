package sendowl

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/sendowl/actions"
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewSendOwl(), Flow, ReadME)

type SendOwl struct{}

func (s *SendOwl) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (s *SendOwl) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		// triggers.NewOrderCompletedTrigger(),
		// triggers.NewProductUpdatedTrigger(),
	}
}

func (s *SendOwl) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListProductsAction(),
		// actions.NewGetProductAction(),
		// actions.NewListOrdersAction(),
		// actions.NewGetOrderAction(),
	}
}

func NewSendOwl() sdk.Integration {
	return &SendOwl{}
}
