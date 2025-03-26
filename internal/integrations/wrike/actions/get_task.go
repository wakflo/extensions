package actions

import (
	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getTaskActionProps struct {
	TaskID string `json:"taskId"`
}

type GetTaskAction struct{}

func (a *GetTaskAction) Name() string {
	return "Get Task"
}

func (a *GetTaskAction) Description() string {
	return "Retrieves detailed information about a specific task in Wrike by its ID."
}

func (a *GetTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getTaskDocs,
	}
}

func (a *GetTaskAction) Icon() *string {
	icon := "mdi:clipboard-text-outline"
	return &icon
}

func (a *GetTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"taskId": shared.GetTaskInput(),
	}
}

func (a *GetTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/tasks/" + input.TaskID

	task, err := shared.GetWrikeClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (a *GetTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetTaskAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *GetTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetTaskAction() sdk.Action {
	return &GetTaskAction{}
}
