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

type createProjectActionProps struct {
	shared.CreateProject
}

type CreateProjectAction struct{}

func (a *CreateProjectAction) Name() string {
	return "Create Project"
}

func (a *CreateProjectAction) Description() string {
	return "Create Project: Initiates the creation of a new project in your project management system, allowing you to start tracking tasks, milestones, and team progress from within our workflow automation platform."
}

func (a *CreateProjectAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateProjectAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createProjectDocs,
	}
}

func (a *CreateProjectAction) Icon() *string {
	return nil
}

func (a *CreateProjectAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("Name").
			SetDescription("Name of the project.").
			SetRequired(true).Build(),
		"project_id": shared.GetProjectsInput(),
		"color": autoform.NewShortTextField().
			SetDisplayName("Color").
			SetDescription("The color of the project icon. Refer to the name column in the Colors guide for more info.").
			SetRequired(false).Build(),
		"is_favorite": autoform.NewBooleanField().
			SetDisplayName("Is Favourite").
			SetDescription("Whether the project is a favorite (a true or false value).").
			SetRequired(false).Build(),
		"view_style": autoform.NewSelectField().
			SetDisplayName("View Style").
			SetDescription("A string value (either list or board, default is list). This determines the way the project is displayed within the Todoist clients.").
			SetOptions(shared.ViewStyleOptions).
			SetRequired(false).Build(),
	}
}

func (a *CreateProjectAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createProjectActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.POST("/projects").Body().AsJSON(input).Send()
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

func (a *CreateProjectAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateProjectAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateProjectAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateProjectAction() sdk.Action {
	return &CreateProjectAction{}
}
