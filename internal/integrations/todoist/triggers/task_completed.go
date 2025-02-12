package triggers

import (
	"context"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type taskCompletedTriggerProps struct {
	ProjectID *string `json:"projectId"`
}

type TaskCompletedTrigger struct{}

func (t *TaskCompletedTrigger) Name() string {
	return "Task Completed"
}

func (t *TaskCompletedTrigger) Description() string {
	return "Triggered when a task is marked as completed in the workflow, allowing subsequent actions to be executed based on the task's status."
}

func (t *TaskCompletedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TaskCompletedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &taskCompletedDocs,
	}
}

func (t *TaskCompletedTrigger) Icon() *string {
	return nil
}

func (t *TaskCompletedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId": shared.GetProjectsInput(),
	}
}

// Start initializes the taskCompletedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *TaskCompletedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the taskCompletedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *TaskCompletedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of taskCompletedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *TaskCompletedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[taskCompletedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Triggered by %v!", input.ProjectID),
	}

	return out, nil
}

func (t *TaskCompletedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *TaskCompletedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TaskCompletedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"content":      "Buy Milk",
		"meta_data":    nil,
		"user_id":      "2671355",
		"task_id":      "2995104339",
		"note_count":   0,
		"project_id":   "2203306141",
		"section_id":   "7025",
		"completed_at": "2015-02-17T15:40:41.000000Z",
		"id":           "1899066186",
	}
}

func NewTaskCompletedTrigger() sdk.Trigger {
	return &TaskCompletedTrigger{}
}
