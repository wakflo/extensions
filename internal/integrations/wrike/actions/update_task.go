package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateTaskActionProps struct {
	TaskID      string `json:"taskId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Importance  string `json:"importance"`
}

type UpdateTaskAction struct{}

func (a *UpdateTaskAction) Name() string {
	return "Update Task"
}

func (a *UpdateTaskAction) Description() string {
	return "Update an existing task in Wrike with new properties such as status."
}

func (a *UpdateTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateTaskDocs,
	}
}

func (a *UpdateTaskAction) Icon() *string {
	icon := "mdi:clipboard-edit-outline"
	return &icon
}

func (a *UpdateTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"taskId": shared.GetTaskInput(),
		"title": autoform.NewShortTextField().
			SetDisplayName("Title").
			SetDescription("The new title of the task.").
			SetRequired(false).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("The new detailed description of the task.").
			SetRequired(false).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetDescription("The new status of the task.").
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
			}).
			Build(),
		"importance": autoform.NewSelectField().
			SetDisplayName("Importance").
			SetDescription("The new importance level of the task.").
			SetRequired(false).
			SetOptions([]*sdkcore.AutoFormSchema{
				{
					Title: "High",
					Const: "High",
				},
				{
					Title: "Normal",
					Const: "Normal",
				},
				{
					Title: "Low",
					Const: "Low",
				},
			}).
			Build(),
	}
}

func (a *UpdateTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Prepare the request data
	data := map[string]interface{}{}

	if input.Title != "" {
		data["title"] = input.Title
	}

	if input.Description != "" {
		data["description"] = input.Description
	}

	if input.Status != "" {
		data["status"] = input.Status
	}

	if input.Importance != "" {
		data["importance"] = input.Importance
	}

	if len(data) == 0 {
		return nil, errors.New("at least one field must be provided to update the task")
	}

	endpoint := "/tasks/" + input.TaskID
	response, err := shared.PutWrikeClient(ctx.Auth.AccessToken, endpoint, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *UpdateTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateTaskAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"id":               "IEADTSKYA5CKABNW",
		"accountId":        "IEADTSKY",
		"title":            "Updated Task",
		"description":      "This task has been updated via the API",
		"briefDescription": "Updated task",
		"parentIds":        []string{"IEADTSKYA5CKAARW"},
		"superParentIds":   []string{"IEADTSKYA5CKAARW"},
		"scope":            "WsFolder",
		"status":           "Completed",
		"importance":       "High",
		"createdDate":      "2023-03-20T14:30:45.000Z",
		"updatedDate":      "2023-03-21T10:15:22.000Z",
		"completedDate":    "2023-03-21T10:15:22.000Z",
		"dates": map[string]interface{}{
			"type":  "Planned",
			"start": "2023-03-21T09:00:00.000Z",
			"due":   "2023-03-25T18:00:00.000Z",
		},
		"responsibleIds": []string{"KUAIJTSKJA", "KUAIJTSKLM"},
		"authorIds":      []string{"KUAIJTSKJA"},
	}
}

func (a *UpdateTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateTaskAction() sdk.Action {
	return &UpdateTaskAction{}
}
