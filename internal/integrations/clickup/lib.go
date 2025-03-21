package clickup

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/clickup/actions"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/extensions/internal/integrations/clickup/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewClickUp(), Flow, ReadME)

type ClickUp struct{}

func (n *ClickUp) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.ClickupSharedAuth,
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
		actions.NewCreateTaskAction(),
		actions.NewCreateFolderOperation(),
		actions.NewCreateFolderlessListOperation(),
		actions.NewCreateSpaceOperation(),
		actions.NewDeleteTaskOperation(),
		actions.NewGetFolderOperation(),
		actions.NewGetFoldersOperation(),
		actions.NewGetFolderlesslistOperation(),
		actions.NewGetListOperation(),
		actions.NewGetSpaceOperation(),
		actions.NewGetSpacesOperation(),
		actions.NewGetTaskOperation(),
		actions.NewGetTasksOperation(),
		actions.NewSearchTaskOperation(),
		actions.NewUpdateSpaceOperation(),
		actions.NewUpdateTaskOperation(),
	}
}

func NewClickUp() sdk.Integration {
	return &ClickUp{}
}
