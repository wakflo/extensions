package wrike

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/wrike/actions"
	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewWrike())

type Wrike struct{}

func (n *Wrike) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (w *Wrike) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.WrikeSharedAuth,
	}
}

func (w *Wrike) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		// triggers.NewTaskCreatedTrigger(),
		// triggers.NewTaskUpdatedTrigger(),
	}
}

func (w *Wrike) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetTaskAction(),
		actions.NewListTasksAction(),
		actions.NewCreateTaskAction(),
		actions.NewUpdateTaskAction(),
	}
}

func NewWrike() sdk.Integration {
	return &Wrike{}
}
