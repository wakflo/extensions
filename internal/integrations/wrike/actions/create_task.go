package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createTaskActionProps struct {
	FolderID    string   `json:"folderId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Importance  string   `json:"importance"`
	StartDate   string   `json:"startDate"`
	DueDate     string   `json:"dueDate"`
	Assignees   []string `json:"assignees"`
}

type CreateTaskAction struct{}

func (a *CreateTaskAction) Name() string {
	return "Create Task"
}

func (a *CreateTaskAction) Description() string {
	return "Create a new task in Wrike with specified title, description, and status."
}

func (a *CreateTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createTaskDocs,
	}
}

func (a *CreateTaskAction) Icon() *string {
	icon := "mdi:clipboard-plus-outline"
	return &icon
}

func (a *CreateTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderId": shared.GetFoldersInput(),
		"title": autoform.NewShortTextField().
			SetDisplayName("Title").
			SetDescription("The title of the task.").
			SetRequired(true).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("The detailed description of the task.").
			SetRequired(false).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetDescription("The status of the task.").
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
			SetDefaultValue("Active").
			Build(),
		"importance": autoform.NewSelectField().
			SetDisplayName("Importance").
			SetDescription("The importance level of the task.").
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
			SetDefaultValue("Normal").
			Build(),
	}
}

func (a *CreateTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"title": input.Title,
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

	endpoint := "/folders/" + input.FolderID + "/tasks"
	response, err := shared.PostWrikeClient(ctx.Auth.AccessToken, endpoint, data)
	if err != nil {
		return nil, err
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format from Wrike API")
	}

	responseData, ok := responseMap["data"].([]interface{})
	if !ok || len(responseData) == 0 {
		return nil, errors.New("invalid response format from Wrike API")
	}

	return responseData[0].(map[string]interface{}), nil
}

func (a *CreateTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateTaskAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"id":               "IEADTSKYA5CKABNW",
		"accountId":        "IEADTSKY",
		"title":            "New Task",
		"description":      "This is a new task created via the API",
		"briefDescription": "New task",
		"parentIds":        []string{"IEADTSKYA5CKAARW"},
		"superParentIds":   []string{"IEADTSKYA5CKAARW"},
		"scope":            "WsFolder",
		"status":           "Active",
		"importance":       "Normal",
		"createdDate":      "2023-03-20T14:30:45.000Z",
		"updatedDate":      "2023-03-20T14:30:45.000Z",
		"dates": map[string]interface{}{
			"type":  "Planned",
			"start": "2023-03-21T09:00:00.000Z",
			"due":   "2023-03-25T18:00:00.000Z",
		},
		"responsibleIds": []string{"KUAIJTSKJA"},
		"authorIds":      []string{"KUAIJTSKJA"},
	}
}

func (a *CreateTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateTaskAction() sdk.Action {
	return &CreateTaskAction{}
}
