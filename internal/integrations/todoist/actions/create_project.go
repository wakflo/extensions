package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/todoist/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createProjectActionProps struct {
	shared.CreateProject
}

type CreateProjectAction struct{}

func (a *CreateProjectAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_project",
		DisplayName:   "Create Project",
		Description:   "Create Project: Initiates the creation of a new project in your project management system, allowing you to start tracking tasks, milestones, and team progress from within our workflow automation platform.",
		Type:          core.ActionTypeAction,
		Documentation: createProjectDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateProjectAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_project", "Create Project")

	form.TextField("name", "Name").
		Placeholder("Name of the project.").
		Required(true).
		HelpText("Name of the project.")

	shared.RegisterProjectsProps(form)

	form.TextField("color", "Color").
		Placeholder("The color of the project icon. Refer to the name column in the Colors guide for more info.").
		HelpText("The color of the project icon. Refer to the name column in the Colors guide for more info.")

	form.CheckboxField("is_favorite", "Is Favourite").
		HelpText("Whether the project is a favorite .")

	form.SelectField("view_style", "View Style").
		AddOption("list", "List").
		AddOption("board", "Board").
		Placeholder("A string value (either list or board, default is list). This determines the way the project is displayed within the Todoist clients.").
		HelpText("A string value (either list or board, default is list). This determines the way the project is displayed within the Todoist clients.")
	schema := form.Build()
	return schema
}

func (a *CreateProjectAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createProjectActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, shared.BaseAPI+"/projects", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authCtx.Token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates an error
	if resp.StatusCode >= 400 {
		// Read the response body to get the error message
		respBody, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(respBody))
	}

	// Read and parse the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var project shared.Project
	err = json.Unmarshal(respBody, &project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (a *CreateProjectAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateProjectAction() sdk.Action {
	return &CreateProjectAction{}
}
