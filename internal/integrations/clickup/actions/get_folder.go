package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

// GetFolderOperation structure and methods
type getFolderProps struct {
	FolderID string `json:"folder-id"`
}

type GetFolderOperation struct{}

// Metadata returns metadata about the action
func (o *GetFolderOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_folder",
		DisplayName:   "Get Folder",
		Description:   "Retrieves details of a specific ClickUp folder by ID.",
		Type:          core.ActionTypeAction,
		Documentation: getFolderDocs,
		Icon:          "material-symbols:folder",
		SampleOutput: map[string]any{
			"id":                "folder123",
			"name":              "Example Folder",
			"orderindex":        "1",
			"override_statuses": false,
			"hidden":            false,
			"space": map[string]any{
				"id":   "space123",
				"name": "Space Name",
			},
			"task_count": "15",
			"lists": []map[string]any{
				{
					"id":      "list123",
					"name":    "List 1",
					"content": "List description",
				},
				{
					"id":      "list456",
					"name":    "List 2",
					"content": "Another list description",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *GetFolderOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_folder", "Get Folder")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	shared.RegisterFoldersInput(form, "Folders", "select a folder", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *GetFolderOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *GetFolderOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getFolderProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	url := "/v2/folder/" + input.FolderID

	folder, _ := shared.GetData(accessToken, url)

	return folder, nil
}

func NewGetFolderOperation() sdk.Action {
	return &GetFolderOperation{}
}
