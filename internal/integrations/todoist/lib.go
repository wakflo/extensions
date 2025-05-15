package todoist

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/todoist/actions"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/extensions/internal/integrations/todoist/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewTodoist())

type Todoist struct{}

func (n *Todoist) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Todoist) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedTodoistAuth,
	}
}

func (n *Todoist) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewTaskCompletedTrigger(),
	}
}

func (n *Todoist) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUpdateTaskAction(),

		actions.NewUpdateProjectAction(),

		actions.NewListTaskAction(),

		actions.NewListProjectCollaboratorsAction(),

		actions.NewListProjectsAction(),

		actions.NewGetProjectAction(),

		actions.NewGetActiveTaskAction(),

		actions.NewCreateTaskAction(),

		actions.NewCreateProjectAction(),
	}
}

func NewTodoist() sdk.Integration {
	return &Todoist{}
}
