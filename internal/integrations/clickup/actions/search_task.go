package actions

import (
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type searchTaskProps struct {
	WorkspaceID string `json:"workspace-id"`
	Query       string `json:"query"`
}

type SearchTaskOperation struct{}

func (o *SearchTaskOperation) Name() string {
	return "Search Task"
}

func (o *SearchTaskOperation) Description() string {
	return "Searches for tasks across a ClickUp workspace based on a query."
}

func (o *SearchTaskOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *SearchTaskOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &searchTaskDocs,
	}
}

func (o *SearchTaskOperation) Icon() *string {
	icon := "material-symbols:search"
	return &icon
}

func (o *SearchTaskOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"query": autoform.NewShortTextField().
			SetDisplayName("Search Query").
			SetDescription("The search query to find tasks").
			SetRequired(true).
			Build(),
	}
}

func (o *SearchTaskOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[searchTaskProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	reqURL := "/v2/team/" + input.WorkspaceID + "/task"
	task, _ := shared.SearchTask(accessToken, reqURL, input.Query)

	return task, nil
}

func (o *SearchTaskOperation) Auth() *sdk.Auth {
	return nil
}

func (o *SearchTaskOperation) SampleData() sdkcore.JSON {
	return map[string]any{
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
				"status": map[string]string{
					"status": "In Progress",
					"color":  "#4286f4",
				},
			},
		},
	}
}

func (o *SearchTaskOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSearchTaskOperation() sdk.Action {
	return &SearchTaskOperation{}
}
