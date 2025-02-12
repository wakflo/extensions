package actions

import (
	"encoding/json"
	"errors"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listTaskActionProps struct {
	Name string `json:"name"`
}

type ListTaskAction struct{}

func (a *ListTaskAction) Name() string {
	return "List Task"
}

func (a *ListTaskAction) Description() string {
	return "List Tasks action retrieves and displays a list of tasks from a specified workflow or project, allowing you to easily view and manage your task portfolio within the automation workflow."
}

func (a *ListTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listTaskDocs,
	}
}

func (a *ListTaskAction) Icon() *string {
	return nil
}

func (a *ListTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"project_id": shared.GetProjectsInput(),
		"section_id": shared.GetSectionsInput(),
		"label": autoform.NewShortTextField().
			SetDisplayName("Label").
			SetDescription("Filter tasks by label name.").
			SetRequired(false).Build(),
		"filter": autoform.NewShortTextField().
			SetDisplayName("Filter").
			SetDescription("Filter by any supported filter. Multiple filters (using the comma , operator) are not supported.").
			SetRequired(false).Build(),
		"lang": autoform.NewShortTextField().
			SetDisplayName("Lang").
			SetDescription("IETF language tag defining what language filter is written in, if differs from default English.").
			SetRequired(false).Build(),
		"ids": autoform.NewArrayField().
			SetDisplayName("IDsss").
			SetDescription("A list of the task IDs to retrieve, this should be a comma separated list.").
			SetItems(
				autoform.NewShortTextField().
					SetDisplayName("ID").
					SetDescription("id").
					SetRequired(true).
					Build(),
			).
			SetRequired(false).Build(),
	}
}

func (a *ListTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	// input, err := integration.InputToTypeSafely[listTaskActionProps](ctx.BaseContext)
	// if err != nil {
	// 	return nil, err
	// }

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.GET("/tasks").Send()
	if err != nil {
		return nil, err
	}

	if rsp.Status().IsError() {
		return nil, errors.New(rsp.Status().Text())
	}

	bytes, err := io.ReadAll(rsp.Raw().Body)
	if err != nil {
		return nil, err
	}

	var tasks []shared.Task
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (a *ListTaskAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListTaskAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListTaskAction() sdk.Action {
	return &ListTaskAction{}
}
