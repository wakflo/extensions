package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type updateTaskActionProps struct {
	TaskID      string `json:"taskId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Importance  string `json:"importance"`
}

type UpdateTaskAction struct{}

func (a *UpdateTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_task",
		DisplayName:   "Update Task",
		Description:   "Update an existing task in Wrike with new properties such as status.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: updateTaskDocs,
		SampleOutput: map[string]interface{}{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *UpdateTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_task", "Update Task")

	shared.GetTaskProp(form)

	form.TextField("title", "Title").
		Placeholder("Enter title").
		Required(false).
		HelpText("Enter title")

	form.TextareaField("description", "Description").
		Placeholder("Enter description").
		Required(false).
		HelpText("The new detailed description of the task.")

	form.SelectField("status", "Status").
		Placeholder("Select status").
		Required(false).
		AddOption("Active", "Active").
		AddOption("Completed", "Completed").
		AddOption("Deferred", "Deferred").
		AddOption("Cancelled", "Cancelled").
		HelpText("Select status")

	form.SelectField("importance", "Importance").
		Placeholder("Select importance").
		Required(false).
		AddOption("High", "High").
		AddOption("Normal", "Normal").
		AddOption("Low", "Low").
		HelpText("The new importance level of the task.")

	schema := form.Build()

	return schema

}

func (a *UpdateTaskAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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
	response, err := shared.PutWrikeClient(token, endpoint, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *UpdateTaskAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewUpdateTaskAction() sdk.Action {
	return &UpdateTaskAction{}
}
