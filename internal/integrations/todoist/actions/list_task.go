package actions

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listTaskActionProps struct {
	Name string `json:"name"`
}

type ListTaskAction struct{}

func (a *ListTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_task",
		DisplayName:   "List Task",
		Description:   "List Tasks action retrieves and displays a list of tasks from a specified workflow or project, allowing you to easily view and manage your task portfolio within the automation workflow.",
		Type:          core.ActionTypeAction,
		Documentation: listTaskDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_task", "List Task")

	shared.RegisterProjectsProps(form)
	shared.RegisterSectionsProps(form)

	form.TextField("label", "Label").
		Placeholder("Filter tasks by label name.").
		HelpText("Filter tasks by label name.")

	form.TextField("filter", "Filter").
		Placeholder("Filter by any supported filter. Multiple filters (using the comma , operator) are not supported.").
		HelpText("Filter by any supported filter. Multiple filters (using the comma , operator) are not supported.")

	form.TextField("lang", "Lang").
		Placeholder("IETF language tag defining what language filter is written in, if differs from default English.").
		HelpText("IETF language tag defining what language filter is written in, if differs from default English.")

	idsArray := form.ArrayField("ids", "IDs")
	idGroup := idsArray.ObjectTemplate("id", "")
	idGroup.TextField("value", "ID").
		Placeholder("Task ID").
		Required(true).
		HelpText("Task ID")

	schema := form.Build()
	return schema
}

func (a *ListTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Comment preserved from original code
	// input, err := integration.InputToTypeSafely[listTaskActionProps](ctx.BaseContext)
	// if err != nil {
	// 	return nil, err
	// }

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(authCtx.AccessToken).
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

func (a *ListTaskAction) Auth() *core.AuthMetadata {
	return nil
}

func NewListTaskAction() sdk.Action {
	return &ListTaskAction{}
}
