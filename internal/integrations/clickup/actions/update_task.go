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

// Metadata returns metadata about the action
func (o *UpdateTaskOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_task",
		DisplayName:   "Update Task",
		Description:   "Updates an existing task in ClickUp with modified details.",
		Type:          core.ActionTypeAction,
		Documentation: updateTaskDocs,
		Icon:          "material-symbols:edit-document",
		SampleOutput: map[string]any{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *UpdateTaskOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_task", "Update Task")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	shared.RegisterFoldersInput(form, "Folders", "select a folder", true)

	shared.RegisterListsInput(form, "Lists", "select a list to create task in", true)

	shared.RegisterTasksInput(form, "Tasks", "select a task to update", true)

	shared.GetAssigneeInput(form, "Assignees", "select a assignee", true)

	form.TextField("name", "name").
		Placeholder("Task Name").
		HelpText("The name of task to update").
		Required(false)

	form.TextareaField("description", "description").
		Placeholder("Task Description").
		HelpText("The description of task to update").
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
func (o *UpdateTaskOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *UpdateTaskOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTaskProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.AccessToken
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

	var response core.JSON
	fmt.Println(response)
	return map[string]interface{}{
		"Report": "Task updated successfully",
	}, nil
}

func NewUpdateTaskOperation() sdk.Action {
	return &UpdateTaskOperation{}
}
