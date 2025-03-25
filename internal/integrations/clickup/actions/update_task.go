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

type updateTaskProps struct {
	TaskID      string `json:"task-id"`
	ListID      string `json:"list-id"`
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
	FolderID    string `json:"folder-id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee-id"`
}

type UpdateTaskOperation struct{}

func (o *UpdateTaskOperation) Name() string {
	return "Update Task"
}

func (o *UpdateTaskOperation) Description() string {
	return "Updates an existing task in ClickUp with modified details."
}

func (o *UpdateTaskOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *UpdateTaskOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateTaskDocs,
	}
}

func (o *UpdateTaskOperation) Icon() *string {
	icon := "material-symbols:edit-document"
	return &icon
}

func (o *UpdateTaskOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace-id": shared.GetWorkSpaceInput("Workspaces", "select a workspace", true),
		"space-id":     shared.GetSpacesInput("Spaces", "select a space", true),
		"folder-id":    shared.GetFoldersInput("Folders", "select a folder", true),
		"list-id":      shared.GetListsInput("Lists", "select a list to create task in", true),
		"task-id":      shared.GetTasksInput("Tasks", "select a task to update", true),
		"assignee-id":  shared.GetAssigneeInput("Assignees", "select a assignee", true),
		"name": autoform.NewShortTextField().
			SetDisplayName("Task Name").
			SetDescription("The name of task to update").
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Task Description").
			SetDescription("The description of task to update").
			Build(),
		"priority": autoform.NewSelectField().
			SetDisplayName("Priority").
			SetDescription("The priority level of the task").
			SetOptions(shared.ClickupPriorityType).
			Build(),
	}
}

func (o *UpdateTaskOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[updateTaskProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	updatedTaskData := map[string]interface{}{}

	if input.Name != "" {
		updatedTaskData["name"] = input.Name
	}
	if input.Description != "" {
		updatedTaskData["description"] = input.Name
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}
		updatedTaskData["priority"] = priority
	}

	if input.AssigneeID != "" {
		assigneeStrings := strings.Split(input.AssigneeID, ",")
		assignees := make([]int, len(assigneeStrings))

		for i, assigneeStr := range assigneeStrings {
			assignee, err := strconv.Atoi(strings.TrimSpace(assigneeStr))
			if err != nil {
				return nil, err
			}
			assignees[i] = assignee
		}
		assigneesObject := map[string][]int{
			"add": assignees,
			"rem": {},
		}

		updatedTaskData["assignees"] = assigneesObject
	}

	taskJSON, err := json.Marshal(updatedTaskData)
	if err != nil {
		return nil, err
	}

	reqURL := shared.BaseURL + "/v2/task/" + input.TaskID
	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(taskJSON))
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
	fmt.Println(response)
	return map[string]interface{}{
		"Report": "Task updated successfully",
	}, nil
}

func (o *UpdateTaskOperation) Auth() *sdk.Auth {
	return nil
}

func (o *UpdateTaskOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":          "abc123",
		"name":        "Updated Task",
		"description": "This task has been updated",
		"status": map[string]string{
			"status": "In Progress",
			"color":  "#4286f4",
		},
		"priority": map[string]any{
			"priority": "Medium",
			"color":    "#f59f00",
		},
		"date_updated": "1647354999999",
	}
}

func (o *UpdateTaskOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateTaskOperation() sdk.Action {
	return &UpdateTaskOperation{}
}
