package prisync

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/prisync/actions"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewPrisync())

type Prisync struct{}

func (n *Prisync) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Prisync) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.PrisyncSharedAuth,
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
