package triggers

import (
	"context"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type taskCompletedTriggerProps struct {
	ProjectID *string `json:"project_id"`
}

type TaskCompletedTrigger struct{}

func (t *TaskCompletedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "task_completed",
		DisplayName:   "Task Completed",
		Description:   "Triggered when a task is marked as completed in the workflow, allowing subsequent actions to be executed based on the task's status.",
		Type:          core.TriggerTypePolling,
		Documentation: taskCompletedDocs,
		SampleOutput: map[string]any{
			"content":      "Buy Milk",
			"meta_data":    nil,
			"user_id":      "2671355",
			"task_id":      "2995104339",
			"note_count":   0,
			"project_id":   "2203306141",
			"section_id":   "7025",
			"completed_at": "2015-02-17T15:40:41.000000Z",
			"id":           "1899066186",
		},
	}
}

func (t *TaskCompletedTrigger) Auth() *core.AuthMetadata {
	return nil
}

func (t *TaskCompletedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("task-completed", "Task Completed")

	shared.RegisterProjectsProps(form)

	schema := form.Build()
	return schema
}

// Start initializes the taskCompletedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *TaskCompletedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the taskCompletedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *TaskCompletedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of taskCompletedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Polling triggers
func (t *TaskCompletedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskCompletedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Triggered by %v!", input.ProjectID),
	}

	return out, nil
}

func (t *TaskCompletedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewTaskCompletedTrigger() sdk.Trigger {
	return &TaskCompletedTrigger{}
}
