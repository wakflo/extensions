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

type getProjectActionProps struct {
	ProjectID string `json:"id"`
}

type GetProjectAction struct{}

func (a *GetProjectAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_project",
		DisplayName:   "Get Project",
		Description:   "Retrieves project details from the specified project management system or platform.",
		Type:          core.ActionTypeAction,
		Documentation: getProjectDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetProjectAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_project", "Get Project")

	form.TextField("id", "Project ID").
		Placeholder("ID of the project.").
		Required(true).
		HelpText("ID of the project.")

	schema := form.Build()
	return schema
}

func (a *GetProjectAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProjectActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(authCtx.Token.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.GET("/projects/" + input.ProjectID).Send()
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

	var project shared.Project
	err = json.Unmarshal(bytes, &project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (a *GetProjectAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetProjectAction() sdk.Action {
	return &GetProjectAction{}
}
