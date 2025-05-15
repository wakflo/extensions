package trello

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/trello/actions"
	"github.com/wakflo/extensions/internal/integrations/trello/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewTrello())

type Trello struct{}

func (n *Trello) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Trello) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.TrelloSharedAuth,
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
