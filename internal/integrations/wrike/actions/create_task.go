package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/wrike/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *CreateTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_task",
		DisplayName:   "Create Task",
		Description:   "Create a new task in Wrike with specified title, description, and status.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createTaskDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_task", "Create Task")

	shared.GetFoldersProp(form)

	form.TextField("title", "Title").
		Placeholder("Enter title").
		Required(true).
		HelpText("Enter title")

	form.TextField("description", "Description").
		Placeholder("Enter description").
		Required(false).
		HelpText("Enter description")

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
		HelpText("Select importance")

	schema := form.Build()

	return schema
}

func (a *CreateTaskAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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
	response, err := shared.PostWrikeClient(token, endpoint, data)
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

func (a *CreateTaskAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateTaskAction() sdk.Action {
	return &CreateTaskAction{}
}
