package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getTaskActionProps struct {
	TaskID string `json:"taskId"`
}

type GetTaskAction struct{}

func (a *GetTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_task",
		DisplayName:   "Get Task",
		Description:   "Retrieves detailed information about a specific task in Wrike by its ID.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getTaskDocs,
		SampleOutput: map[string]any{
			"id":               "IEADTSKYA5CKABJN",
			"accountId":        "IEADTSKY",
			"title":            "Example Task",
			"description":      "This is an example task",
			"briefDescription": "Example task",
			"parentIds":        []string{"IEADTSKYA5CKAARW"},
			"superParentIds":   []string{"IEADTSKYA5CKAARW"},
			"scope":            "WsFolder",
			"status":           "Active",
			"importance":       "Normal",
			"createdDate":      "2023-03-15T10:30:45.000Z",
			"updatedDate":      "2023-03-15T15:20:15.000Z",
			"completedDate":    nil,
			"dates": map[string]any{
				"type":     "Planned",
				"duration": "86400",
				"start":    "2023-03-16T09:00:00.000Z",
				"due":      "2023-03-17T18:00:00.000Z",
			},
			"customFields":   []map[string]any{},
			"customStatus":   nil,
			"responsibleIds": []string{"KUAIJTSKJA"},
			"authorIds":      []string{"KUAIJTSKJA"},
			"hasAttachments": false,
			"priority":       "Normal",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_task", "Get Task")

	shared.GetTaskProp(form)

	schema := form.Build()

	return schema
}

func (a *GetTaskAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	url := "/tasks/" + input.TaskID

	task, err := shared.GetWrikeClient(token, url)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (a *GetTaskAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetTaskAction() sdk.Action {
	return &GetTaskAction{}
}
