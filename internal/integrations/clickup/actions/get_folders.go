package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getFoldersProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
}

type GetFoldersOperation struct{}

// Metadata returns metadata about the action
func (o *GetFoldersOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_folders",
		DisplayName:   "Get Folders",
		Description:   "Retrieves all folders within a specified ClickUp space.",
		Type:          core.ActionTypeAction,
		Documentation: getFoldersDocs,
		Icon:          "material-symbols:folder-copy",
		SampleOutput: map[string]any{
			"folders": []map[string]any{
				{
					"id":                "folder123",
					"name":              "Folder 1",
					"orderindex":        "1",
					"override_statuses": false,
					"hidden":            false,
					"task_count":        "8",
				},
				{
					"id":                "folder456",
					"name":              "Folder 2",
					"orderindex":        "2",
					"override_statuses": false,
					"hidden":            false,
					"task_count":        "12",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *GetFoldersOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_folders", "Get Folders")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *GetFoldersOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *GetFoldersOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getFoldersProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	url := "/v2/space/" + input.SpaceID + "/folder"

	folders, err := shared.GetData(accessToken, url)
	if err != nil {
		return nil, err
	}

	return folders, nil
}

func NewGetFoldersOperation() sdk.Action {
	return &GetFoldersOperation{}
}
