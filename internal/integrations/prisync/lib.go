package prisync

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/prisync/actions"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewPrisync(), Flow, ReadME)

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
