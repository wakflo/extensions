package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listProjectCollaboratorsActionProps struct {
	ProjectID string `json:"project_id"`
}

type ListProjectCollaboratorsAction struct{}

func (a *ListProjectCollaboratorsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_project_collaborators",
		DisplayName:   "List Project Collaborators",
		Description:   "Retrieves a list of users who are currently collaborating on a specific project, including their roles and permissions.",
		Type:          core.ActionTypeAction,
		Documentation: listProjectCollaboratorsDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListProjectCollaboratorsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_project_collaborators", "List Project Collaborators")

	shared.RegisterProjectsProps(form)

	schema := form.Build()
	return schema
}

func (a *ListProjectCollaboratorsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listProjectCollaboratorsActionProps](ctx)
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

	var collaborators []shared.Collaborator
	err = json.Unmarshal(bytes, &collaborators)
	if err != nil {
		return nil, err
	}

	return collaborators, nil
}

func (a *ListProjectCollaboratorsAction) Auth() *core.AuthMetadata {
	return nil
}

func NewListProjectCollaboratorsAction() sdk.Action {
	return &ListProjectCollaboratorsAction{}
}
