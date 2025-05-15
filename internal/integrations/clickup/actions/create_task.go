package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createTaskActionProps struct {
	ListID      string `json:"list-id"`
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
	FolderID    string `json:"folder-id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee-id"`
}

type CreateTaskAction struct{}

// Metadata returns metadata about the action
func (a *CreateTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_task",
		DisplayName:   "Create Task",
		Description:   "Creates a new task in a specified ClickUp list with customizable details like name, description, priority, and assignees.",
		Type:          core.ActionTypeAction,
		Documentation: createTaskDocs,
		Icon:          "material-symbols:add-task",
		SampleOutput: map[string]any{
			"id":          "abc123",
			"name":        "Example Task",
			"description": "This is a sample task",
			"status": map[string]string{
				"status": "Open",
				"color":  "#d3d3d3",
			},
			"priority": map[string]any{
				"priority": "High",
				"color":    "#f50000",
			},
			"date_created": "1647354847362",
			"date_updated": "1647354847362",
			"creator": map[string]any{
				"id":       "123456",
				"username": "John Doe",
				"email":    "john@example.com",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_task", "Create Task")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	shared.RegisterFoldersInput(form, "Folders", "select a folder", true)

	shared.RegisterListsInput(form, "Lists", "select a list to create task in", true)

	shared.GetAssigneeInput(form, "Assignees", "select a assignee", true)

	form.TextField("name", "name").
		Placeholder("Task Name").
		HelpText("The name of the task").
		Required(true)

	form.TextareaField("description", "description").
		Placeholder("Task Description").
		HelpText("The description of the task").
		Required(false)

	form.SelectField("priority", "priority").
		Placeholder("Priority").
		HelpText("The priority level of the task").
		AddOptions(shared.ClickupPriorityType...).
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateTaskAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.AccessToken

	taskData := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}
		taskData["priority"] = priority
	}

	if input.AssigneeID != "" {
		assigneeStrings := strings.Split(input.AssigneeID, ",")
		assignees := make([]int, len(assigneeStrings))

		for i, assigneeStr := range assigneeStrings {
			assignee, err := strconv.Atoi(strings.TrimSpace(assigneeStr))
			if err != nil {
				fmt.Println("Assignee conversion error:", err)
				return nil, err
			}
			assignees[i] = assignee
		}
		taskData["assignees"] = assignees
	}

	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		return nil, err
	}

	reqURL := shared.BaseURL + "/v2/list/" + input.ListID + "/task"
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(taskJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response core.JSON
	if newErr := json.NewDecoder(res.Body).Decode(&response); newErr != nil {
		return nil, err
	}

	return response, nil
}

func NewCreateTaskAction() sdk.Action {
	return &CreateTaskAction{}
}
