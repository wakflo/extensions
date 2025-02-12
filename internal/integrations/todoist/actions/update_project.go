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

type updateProjectActionProps struct {
	shared.UpdateProject
	ProjectID string `json:"id"`
}

type UpdateProjectAction struct{}

func (a *UpdateProjectAction) Name() string {
	return "Update Project"
}

func (a *UpdateProjectAction) Description() string {
	return "Updates project information in the designated project management system, ensuring accurate and up-to-date records of project details, milestones, and tasks."
}

func (a *UpdateProjectAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateProjectAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateProjectDocs,
	}
}

func (a *UpdateProjectAction) Icon() *string {
	return nil
}

func (a *UpdateProjectAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Project ID").
			SetDescription("ID of the project.").
			SetRequired(true).Build(),
		"name": autoform.NewShortTextField().
			SetDisplayName("Name").
			SetDescription("Name of the project.").
			SetRequired(false).Build(),
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

func (a *UpdateProjectAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateProjectActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
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

func (a *UpdateProjectAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateProjectAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateProjectAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateProjectAction() sdk.Action {
	return &UpdateProjectAction{}
}
