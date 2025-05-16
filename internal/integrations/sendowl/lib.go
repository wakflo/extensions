package sendowl

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/sendowl/actions"
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	"github.com/wakflo/extensions/internal/integrations/sendowl/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewSendOwl())

type SendOwl struct{}

func (n *SendOwl) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (s *SendOwl) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
	}
}

func (s *SendOwl) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewOrderCompletedTrigger(),
		triggers.NewProductUpdatedTrigger(),
	}
}

func (s *SendOwl) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListProductsAction(),
		actions.NewGetProductAction(),
		actions.NewListOrdersAction(),
		actions.NewGetOrderAction(),
	}
}

func NewSendOwl() sdk.Integration {
	return &SendOwl{}
}
