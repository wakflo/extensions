package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type searchTaskProps struct {
	WorkspaceID string `json:"workspace-id"`
	Query       string `json:"query"`
}

type SearchTaskOperation struct{}

func (o *SearchTaskOperation) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (o *SearchTaskOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "search_task",
		DisplayName:   "Search Task",
		Description:   "Searches for tasks across a ClickUp workspace based on a query.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: searchTaskDocs,
		SampleOutput: map[string]any{
			"tasks": []map[string]any{
				{
					"id":   "abc123",
					"name": "Found Task 1",
					"status": map[string]string{
						"status": "Open",
						"color":  "#d3d3d3",
					},
				},
				{
					"id":   "def456",
					"name": "Found Task 2",
				},
			},
		},
		Icon: "material-symbols:search",
	}
}

func (o *SearchTaskOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("search-task", "Search for Tasks")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	form.TextareaField("query", "Search Query").
		Placeholder("Search for tasks").
		Required(true).
		HelpText("The search query to find tasks")

	return form.Build()
}

func (o *SearchTaskOperation) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth().AccessToken
	input, err := sdk.InputToTypeSafely[searchTaskProps](ctx)
	if err != nil {
		return nil, err
	}

	reqURL := "/v2/team/" + input.WorkspaceID + "/task"
	task, _ := shared.SearchTask(accessToken, reqURL, input.Query)

	return task, nil
}

func (o *SearchTaskOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSearchTaskAction() sdk.Action {
	return &SearchTaskOperation{}
}
