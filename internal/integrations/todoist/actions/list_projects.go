package actions

import (
	"encoding/json"
	"errors"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listProjectsActionProps struct {
	ProjectID string `json:"project_id"`
}

type ListProjectsAction struct{}

func (a *ListProjectsAction) Name() string {
	return "List Projects"
}

func (a *ListProjectsAction) Description() string {
	return "Retrieves a list of all projects within your organization, allowing you to easily manage and track multiple projects from a single location."
}

func (a *ListProjectsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListProjectsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listProjectsDocs,
	}
}

func (a *ListProjectsAction) Icon() *string {
	return nil
}

func (a *ListProjectsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"project_id": shared.GetProjectsInput(),
	}
}

func (a *ListProjectsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	client := fastshot.NewClient(shared.BaseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
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

func (a *ListProjectsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListProjectsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListProjectsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListProjectsAction() sdk.Action {
	return &ListProjectsAction{}
}
