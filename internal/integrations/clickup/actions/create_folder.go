package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

// CreateFolderOperation structure and methods
type createFolderProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
	Name        string `json:"name"`
}

type CreateFolderOperation struct{}

// Metadata returns metadata about the action
func (o *CreateFolderOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_folder",
		DisplayName:   "Create Folder",
		Description:   "Creates a new folder in a specified ClickUp space.",
		Type:          core.ActionTypeAction,
		Documentation: createFolderDocs,
		Icon:          "material-symbols:create-new-folder",
		SampleOutput: map[string]any{
			"id":                "folder123",
			"name":              "New Folder",
			"orderindex":        "1",
			"override_statuses": false,
			"hidden":            false,
			"space": map[string]any{
				"id":   "space123",
				"name": "Space Name",
			},
			"task_count": "0",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *CreateFolderOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_folder", "Create Folder")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	form.TextField("name", "name").
		Placeholder("Folder Name").
		HelpText("The name of the folder").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *CreateFolderOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *CreateFolderOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createFolderProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	reqURL := "/v2/space/" + input.SpaceID + "/folder"

	folder, err := shared.CreateItem(accessToken, input.Name, reqURL)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func NewCreateFolderOperation() sdk.Action {
	return &CreateFolderOperation{}
}
