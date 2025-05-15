package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getFolderlessListProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
}

type GetFolderlesslistOperation struct{}

// Metadata returns metadata about the action
func (o *GetFolderlesslistOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_folderless_lists",
		DisplayName:   "Get Folderless Lists",
		Description:   "Retrieves all lists directly in a space that do not belong to any folder.",
		Type:          core.ActionTypeAction,
		Documentation: getFolderlessListDocs,
		Icon:          "material-symbols:folder-off-outline",
		SampleOutput: map[string]any{
			"lists": []map[string]any{
				{
					"id":         "list123",
					"name":       "Folderless List 1",
					"content":    "First list",
					"orderindex": "1",
				},
				{
					"id":         "list456",
					"name":       "Folderless List 2",
					"content":    "Second list",
					"orderindex": "2",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *GetFolderlesslistOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_folderless_lists", "Get Folderless Lists")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *GetFolderlesslistOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *GetFolderlesslistOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getFolderlessListProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	folderlessList, _ := shared.GetSpace(accessToken, input.SpaceID)

	return folderlessList, nil
}

func NewGetFolderlesslistOperation() sdk.Action {
	return &GetFolderlesslistOperation{}
}
