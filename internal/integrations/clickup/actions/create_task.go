package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *CreateTaskAction) Name() string {
	return "Create Task"
}

func (a *CreateTaskAction) Description() string {
	return "Creates a new task in a specified ClickUp list with customizable details like name, description, priority, and assignees."
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
	icon := "material-symbols:add-task"
	return &icon
}

func (a *CreateTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
		"folder-id":    shared.GetFoldersInput("Folders", "select a folder", true),
		"list-id":      shared.GetListsInput("Lists", "select a list to create task in", true),
		"assignee-id":  shared.GetAssigneeInput("Assignees", "select a assignee", true),
		"name": autoform.NewShortTextField().
			SetDisplayName("Task Name").
			SetDescription("The name of the task").
			SetRequired(true).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Task Description").
			SetDescription("The description of the task").
			Build(),
		"priority": autoform.NewSelectField().
			SetDisplayName("Priority").
			SetDescription("The priority level of the task").
			SetOptions(shared.ClickupPriorityType).
			Build(),
	}
}

func (a *CreateTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	accessToken := ctx.Auth.AccessToken

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

	var response sdk.JSON
	if newErr := json.NewDecoder(res.Body).Decode(&response); newErr != nil {
		return nil, err
	}

	return response, nil
}

func (a *CreateTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateTaskAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *CreateTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateTaskAction() sdk.Action {
	return &CreateTaskAction{}
}
