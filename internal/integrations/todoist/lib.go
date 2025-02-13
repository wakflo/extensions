package todoist

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/todoist/actions"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/extensions/internal/integrations/todoist/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewTodoist(), Flow, ReadME)

type Todoist struct{}

func (n *Todoist) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
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
