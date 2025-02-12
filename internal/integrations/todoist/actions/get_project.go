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

type getProjectActionProps struct {
	ProjectID string `json:"id"`
}

type GetProjectAction struct{}

func (a *GetProjectAction) Name() string {
	return "Get Project"
}

func (a *GetProjectAction) Description() string {
	return "Retrieves project details from the specified project management system or platform."
}

func (a *GetProjectAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetProjectAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getProjectDocs,
	}
}

func (a *GetProjectAction) Icon() *string {
	return nil
}

func (a *GetProjectAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Project ID").
			SetDescription("ID of the project.").
			SetRequired(true).Build(),
	}
}

func (a *GetProjectAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProjectActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
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

func (a *GetProjectAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetProjectAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetProjectAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetProjectAction() sdk.Action {
	return &GetProjectAction{}
}
