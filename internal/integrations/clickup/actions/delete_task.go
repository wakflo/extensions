package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type deleteTaskProps struct {
	TaskID string `json:"task-id"`
}

type DeleteTaskOperation struct{}

func (o *DeleteTaskOperation) Name() string {
	return "Delete Task"
}

func (o *DeleteTaskOperation) Description() string {
	return "Deletes a task from ClickUp."
}

func (o *DeleteTaskOperation) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (o *DeleteTaskOperation) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &deleteTaskDocs,
	}
}

func (o *DeleteTaskOperation) Icon() *string {
	icon := "material-symbols:delete"
	return &icon
}

func (o *DeleteTaskOperation) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"task-id": autoform.NewShortTextField().
			SetDisplayName("Task ID").
			SetDescription("The ID of the task to delete").
			SetRequired(true).
			Build(),
	}
}

func (o *DeleteTaskOperation) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	accessToken := ctx.Auth.AccessToken
	input, err := sdk.InputToTypeSafely[deleteTaskProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	reqURL := shared.BaseURL + "/v2/task/" + input.TaskID
	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return map[string]interface{}{
		"task": "Task Deleted",
	}, nil
}

func (o *DeleteTaskOperation) Auth() *sdk.Auth {
	return nil
}

func (o *DeleteTaskOperation) SampleData() sdkcore.JSON {
	return map[string]any{
		"success": true,
	}
}

func (o *DeleteTaskOperation) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDeleteTaskOperation() sdk.Action {
	return &DeleteTaskOperation{}
}
