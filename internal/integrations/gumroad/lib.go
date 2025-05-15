package gumroad

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/gumroad/actions"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGumroad())

type Gumroad struct{}

func (n *Gumroad) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Gumroad) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedGumroadAuth,
	}
}

func (n *Gumroad) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Gumroad) Actions() []sdk.Action {
	return []sdk.Action{
		// actions.NewGetProductAction(),
		actions.NewListProductsAction(),
		actions.NewGetProductAction(),
		actions.NewDisableProductAction(),
		actions.NewEnableProductAction(),
		actions.NewListSalesAction(),
		actions.NewGetSaleAction(),
		actions.NewDeleteProductAction(),
		actions.NewMarkasShippedAction(),
	}
}

func NewGumroad() sdk.Integration {
	return &Gumroad{}
}
