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

type updateProjectActionProps struct {
	shared.UpdateProject
	ProjectID string `json:"project_id"`
}

type UpdateProjectAction struct{}

func (a *UpdateProjectAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_project",
		DisplayName:   "Update Project",
		Description:   "Updates project information in the designated project management system, ensuring accurate and up-to-date records of project details, milestones, and tasks.",
		Type:          core.ActionTypeAction,
		Documentation: updateProjectDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *UpdateProjectAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_project", "Update Project")

	form.TextField("id", "Project ID").
		Placeholder("ID of the project.").
		Required(true).
		HelpText("ID of the project.")

	form.TextField("name", "Name").
		Placeholder("Name of the project.").
		HelpText("Name of the project.")

	form.TextField("color", "Color").
		Placeholder("The color of the project icon. Refer to the name column in the Colors guide for more info.").
		HelpText("The color of the project icon. Refer to the name column in the Colors guide for more info.")

	form.CheckboxField("is_favorite", "Is Favourite").
		HelpText("Whether the project is a favorite (a true or false value).")

	form.SelectField("view_style", "View Style").
		AddOption("list", "List").
		AddOption("board", "Board").
		Placeholder("A string value (either list or board, default is list). This determines the way the project is displayed within the Todoist clients.").
		HelpText("A string value (either list or board, default is list). This determines the way the project is displayed within the Todoist clients.")

	schema := form.Build()
	return schema
}

func (a *UpdateProjectAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateProjectActionProps](ctx)
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

	rsp, err := client.POST("/projects/" + input.ProjectID).Body().AsJSON(input.UpdateProject).Send()
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

func (a *UpdateProjectAction) Auth() *core.AuthMetadata {
	return nil
}

func NewUpdateProjectAction() sdk.Action {
	return &UpdateProjectAction{}
}
