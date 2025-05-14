package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/monday/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type createItemActionProps struct {
	WorkspaceID           string `json:"workspace_id,omitempty"`
	BoardID               string `json:"board_id,omitempty"`
	GroupID               string `json:"group_id,omitempty"`
	ItemName              string `json:"item_name"`
	CreateLabelsIfMissing bool   `json:"create_labels_if_missing,omitempty"`
}

type CreateItemAction struct{}

// func (a *CreateItemAction) Name() string {
// 	return "Create Item"
// }

func (a *CreateItemAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_item",
		DisplayName:   "Create Item",
		Description:   "Create Item: Automatically generates a new item in your system, such as a task, issue, or project, with customizable fields and attributes.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createItemDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateItemAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("create_item", "Create Item")

	shared.GetWorkspaceProp(form)

	shared.GetBoardProp("board_id", "Board ID", "Select Board", form)

	form.TextField("item_name", "Item Name").
		Placeholder("New Item").
		Required(true).
		HelpText("Item name.")

	form.CheckboxField("create_labels_if_missing", "Create Labels if Missing").
		HelpText("Creates status/dropdown labels if they are missing. This requires permission to change the board structure.").
		DefaultValue(false).
		Required(false)

	schema := form.Build()

	return schema
}

func (a *CreateItemAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createItemActionProps](ctx)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	fields["item_name"] = fmt.Sprintf(`"%s"`, input.ItemName)
	fields["board_id"] = fmt.Sprintf(`"%s"`, input.BoardID)

	if input.GroupID != "" {
		fields["group_id"] = fmt.Sprintf(`"%s"`, input.GroupID)
	}

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	query := fmt.Sprintf(`
		mutation {
  create_item (%s) {
    	id
		name
  }
}`, strings.Join(fieldStrings, "\n"))

	response, err := shared.MondayClient(ctx, query)
	if err != nil {
		return nil, err
	}

	item, ok := response["data"].(map[string]interface{})["create_item"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract item from response")
	}

	return item, nil
}

func (a *CreateItemAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateItemAction() sdk.Action {
	return &CreateItemAction{}
}
