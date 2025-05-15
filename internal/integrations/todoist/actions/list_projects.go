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

type listProjectsActionProps struct {
	ProjectID string `json:"project_id"`
}

type ListProjectsAction struct{}

func (a *ListProjectsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_projects",
		DisplayName:   "List Projects",
		Description:   "Retrieves a list of all projects within your organization, allowing you to easily manage and track multiple projects from a single location.",
		Type:          core.ActionTypeAction,
		Documentation: listProjectsDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListProjectsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_projects", "List Projects")

	shared.RegisterProjectsProps(form)

	schema := form.Build()
	return schema
}

func (a *ListProjectsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(authCtx.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.GET("/projects").Send()
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

	var projects []shared.Project
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (a *ListProjectsAction) Auth() *core.AuthMetadata {
	return nil
}

func NewListProjectsAction() sdk.Action {
	return &ListProjectsAction{}
}
