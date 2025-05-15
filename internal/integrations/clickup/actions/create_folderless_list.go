package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createFolderlessListProps struct {
	SpaceID     string `json:"space-id"`
	WorkspaceID string `json:"workspace-id"`
	Name        string `json:"name"`
}

type CreateFolderlessListOperation struct{}

// Metadata returns metadata about the action
func (o *CreateFolderlessListOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_folderless_list",
		DisplayName:   "Create Folderless List",
		Description:   "Creates a new list directly in a space without a parent folder.",
		Type:          core.ActionTypeAction,
		Documentation: createFolderlessListDocs,
		Icon:          "material-symbols:add-box-outline",
		SampleOutput: map[string]any{
			"id":         "list123",
			"name":       "New Folderless List",
			"content":    "List description",
			"orderindex": 1,
			"status": map[string]any{
				"status": "Open",
				"color":  "#d3d3d3",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *CreateFolderlessListOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_folderless_list", "Create Folderless List")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "Spaces", "select a space", true)

	form.TextField("name", "name").
		Placeholder("List Name").
		HelpText("The name of the list").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *CreateFolderlessListOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *CreateFolderlessListOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createFolderlessListProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	reqURL := "/v2/space/" + input.SpaceID + "/list"

	response, err := shared.CreateItem(accessToken, input.Name, reqURL)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewCreateFolderlessListOperation() sdk.Action {
	return &CreateFolderlessListOperation{}
}
