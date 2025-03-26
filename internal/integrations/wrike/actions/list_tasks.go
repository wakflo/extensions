package actions

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listTasksActionProps struct {
	FolderID string `json:"folderId"`
	Status   string `json:"status"`
	Limit    int    `json:"limit"`
}

type ListTasksAction struct{}

func (a *ListTasksAction) Name() string {
	return "List Tasks"
}

func (a *ListTasksAction) Description() string {
	return "Retrieve a list of tasks from Wrike with optional filtering parameters."
}

func (a *ListTasksAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListTasksAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listTasksDocs,
	}
}

func (a *ListTasksAction) Icon() *string {
	icon := "mdi:clipboard-list-outline"
	return &icon
}

func (a *ListTasksAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderId": shared.GetFoldersInput(),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetDescription("Filter tasks by their status.").
			SetRequired(false).
			SetOptions([]*sdkcore.AutoFormSchema{
				{
					Title: "Active",
					Const: "Active",
				},
				{
					Title: "Completed",
					Const: "Completed",
				},
				{
					Title: "Deferred",
					Const: "Deferred",
				},
				{
					Title: "Cancelled",
					Const: "Cancelled",
				},
				{
					Title: "All",
					Const: "All",
				},
			}).
			SetDefaultValue("Active").
			Build(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of tasks to return (1-100).").
			SetRequired(false).
			Build(),
	}
}

func (a *ListTasksAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listTasksActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var endpoint string
	var queryParams []string

	if input.FolderID != "" {
		endpoint = fmt.Sprintf("/folders/%s/tasks", input.FolderID)
	} else {
		endpoint = "/tasks"
	}

	if input.Status != "" && input.Status != "All" {
		queryParams = append(queryParams, fmt.Sprintf("status=%s", url.QueryEscape(input.Status)))
	}

	if input.Limit > 0 {
		queryParams = append(queryParams, fmt.Sprintf("limit=%d", input.Limit))
	}

	if len(queryParams) > 0 {
		endpoint = endpoint + "?" + strings.Join(queryParams, "&")
	}

	response, err := shared.GetWrikeClient(ctx.Auth.AccessToken, endpoint)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *ListTasksAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListTasksAction) SampleData() sdkcore.JSON {
	return []interface{}{
		map[string]interface{}{
			"id":          "IEADTSKYA5CKABJN",
			"accountId":   "IEADTSKY",
			"title":       "Example Task 1",
			"description": "This is an example task",
			"status":      "Active",
			"importance":  "Normal",
			"createdDate": "2023-03-15T10:30:45.000Z",
			"updatedDate": "2023-03-15T15:20:15.000Z",
		},
		map[string]interface{}{
			"id":          "IEADTSKYA5CKABJ2",
			"accountId":   "IEADTSKY",
			"title":       "Example Task 2",
			"description": "This is another example task",
			"status":      "Active",
			"importance":  "High",
			"createdDate": "2023-03-16T09:15:30.000Z",
			"updatedDate": "2023-03-16T14:45:22.000Z",
		},
	}
}

func (a *ListTasksAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListTasksAction() sdk.Action {
	return &ListTasksAction{}
}
