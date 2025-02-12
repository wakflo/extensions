package prisync

import (
	"github.com/wakflo/extensions/internal/integrations/prisync/actions"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewPrisync())

type Prisync struct{}

func (n *Prisync) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *Prisync) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Prisync) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewAddBatchProductAction(),

		actions.NewGetProductsSmartpriceAction(),

		actions.NewListProductsAction(),

		actions.NewGetProductAction(),

		actions.NewEditProductAction(),

		actions.NewAddProductAction(),
	}
}

func NewPrisync() sdk.Integration {
	return &Prisync{}
}
