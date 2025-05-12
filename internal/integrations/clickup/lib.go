package clickup

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/clickup/actions"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/extensions/internal/integrations/clickup/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewClickUp())

type ClickUp struct{}

func (n *ClickUp) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *ClickUp) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.ClickupSharedAuth,
	}
}

func (n *ClickUp) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewTaskCreatedTrigger(),
		triggers.NewTaskUpdatedTrigger(),
	}
}

func (n *ClickUp) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetSpacesAction(),
		actions.NewGetTaskAction(),
		actions.NewGetTasksAction(),
	}
}

func NewClickUp() sdk.Integration {
	return &ClickUp{}
}
