package actions

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listTasksActionProps struct {
	FolderID string `json:"folderId"`
	Status   string `json:"status"`
	Limit    int    `json:"limit"`
}

type ListTasksAction struct{}

func (a *ListTasksAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_tasks",
		DisplayName:   "List Tasks",
		Description:   "Retrieve a list of tasks from Wrike with optional filtering parameters.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listTasksDocs,
		SampleOutput: []interface{}{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *ListTasksAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_tasks", "List Tasks")

	shared.GetFoldersProp(form)

	form.SelectField("status", "Status").
		Placeholder("Select status").
		Required(false).
		AddOption("Active", "Active").
		AddOption("Completed", "Completed").
		AddOption("Deferred", "Deferred").
		AddOption("Cancelled", "Cancelled").
		AddOption("All", "All").
		HelpText("Filter tasks by their status.")

	form.NumberField("limit", "Limit").
		Placeholder("Enter limit").
		Required(false).
		HelpText("Maximum number of tasks to return (1-100).")

	schema := form.Build()

	return schema

}

func (a *ListTasksAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listTasksActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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

	response, err := shared.GetWrikeClient(token, endpoint)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *ListTasksAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListTasksAction() sdk.Action {
	return &ListTasksAction{}
}
