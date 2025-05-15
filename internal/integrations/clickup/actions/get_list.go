package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/clickup/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getListProps struct {
	ListID string `json:"list-id"`
}

type GetListOperation struct{}

// Metadata returns metadata about the action
func (o *GetListOperation) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_list",
		DisplayName:   "Get List",
		Description:   "Retrieves details of a specific ClickUp list by ID.",
		Type:          core.ActionTypeAction,
		Documentation: getListDocs,
		Icon:          "material-symbols:format-list-bulleted",
		SampleOutput: map[string]any{
			"id":      "list123",
			"name":    "Example List",
			"content": "List description",
			"statuses": []map[string]any{
				{
					"id":     "st123",
					"status": "Open",
					"color":  "#d3d3d3",
				},
			},
			"task_count": "24",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (o *GetListOperation) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_list", "Get List")

	shared.RegisterWorkSpaceInput(form, "Workspaces", "select a workspace", true)

	shared.RegisterSpacesInput(form, "space-id", "select a space", true)

	shared.RegisterFoldersInput(form, "Folders", "select a folder", true)

	shared.RegisterListsInput(form, "Lists", "select a list to create task in", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (o *GetListOperation) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (o *GetListOperation) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getListProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessToken := authCtx.Token.AccessToken
	list, err := shared.GetList(accessToken, input.ListID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func NewGetListOperation() sdk.Action {
	return &GetListOperation{}
}
