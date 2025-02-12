package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listProjectCollaboratorsActionProps struct {
	ProjectID string `json:"id"`
}

type ListProjectCollaboratorsAction struct{}

func (a *ListProjectCollaboratorsAction) Name() string {
	return "List Project Collaborators"
}

func (a *ListProjectCollaboratorsAction) Description() string {
	return "Retrieves a list of users who are currently collaborating on a specific project, including their roles and permissions."
}

func (a *ListProjectCollaboratorsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListProjectCollaboratorsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listProjectCollaboratorsDocs,
	}
}

func (a *ListProjectCollaboratorsAction) Icon() *string {
	return nil
}

func (a *ListProjectCollaboratorsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Project ID").
			SetDescription("ID of the project.").
			SetRequired(true).Build(),
	}
}

func (a *ListProjectCollaboratorsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listProjectCollaboratorsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.GET(fmt.Sprintf("/projects/%s/collaborators", input.ProjectID)).Send()
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

	var projects []shared.Collaborator
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (a *ListProjectCollaboratorsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListProjectCollaboratorsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListProjectCollaboratorsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListProjectCollaboratorsAction() sdk.Action {
	return &ListProjectCollaboratorsAction{}
}
