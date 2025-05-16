package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type taskCreatedTriggerActionProps struct {
	ListID      string  `json:"list-id"`
	WorkspaceID *string `json:"workspace-id"`
}

type TaskCreatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *TaskCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_created",
		DisplayName:   "Task Created",
		Description:   "Triggers a workflow when a new task is created in a specified ClickUp workspace or list, allowing you to automate subsequent actions based on new task creation events.",
		Type:          core.TriggerTypePolling,
		Documentation: taskCreatedDocs,
		SampleOutput: map[string]any{
			"id":   "abc123",
			"name": "New Task",
			"status": map[string]any{
				"status": "Open",
				"color":  "#d3d3d3",
			},
			"date_created": "1647354847362",
			"creator": map[string]any{
				"id":       "123456",
				"username": "John Doe",
				"email":    "john@example.com",
			},
		},
	}
}

// Props returns the schema for the trigger's input configuration
func (t *TaskCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task_created", "Task Created")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)
	shared.RegisterListsInput(form, "Lists", "select a list to create task in", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the trigger
func (t *TaskCreatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// Start initializes the trigger
func (t *TaskCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger
func (t *TaskCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main trigger logic
func (t *TaskCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[taskCreatedTriggerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.ListID == "" && (input.WorkspaceID == nil || *input.WorkspaceID == "") {
		return nil, errors.New("either List ID or Workspace ID must be provided")
	}

	// Get the last run time
	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	var createdSince string
	if lastRunTime != nil {
		lastRunTimeValue := lastRunTime.(*time.Time)
		createdSince = strconv.FormatInt(lastRunTimeValue.UnixNano()/int64(time.Millisecond), 10)
	}

	var endpoint string
	var queryParam string

	if input.ListID != "" {
		endpoint = "/v2/list/" + input.ListID + "/task"
	} else {
		endpoint = "/v2/team/" + *input.WorkspaceID + "/task"
	}

	if createdSince != "" {
		queryParam = "?date_created_gt=" + createdSince
		endpoint += queryParam
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	response, err := shared.GetClickUpClient(authCtx.Token.AccessToken, endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	tasksArray, ok := response["tasks"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: tasks field is not an array")
	}

	if len(tasksArray) == 0 {
		return []interface{}{}, nil
	}

	return tasksArray, nil
}

// Criteria returns the criteria for the trigger
func (t *TaskCreatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewTaskCreatedTrigger() sdk.Trigger {
	return &TaskCreatedTrigger{}
}
