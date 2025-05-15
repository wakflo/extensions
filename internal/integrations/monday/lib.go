package monday

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/monday/actions"
	"github.com/wakflo/extensions/internal/integrations/monday/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewMonday())

type Monday struct{}

func (n *Monday) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Monday) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
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
