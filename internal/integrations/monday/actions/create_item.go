package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/monday/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createItemActionProps struct {
	WorkspaceID           string `json:"workspace_id,omitempty"`
	BoardID               string `json:"board_id,omitempty"`
	GroupID               string `json:"group_id,omitempty"`
	ItemName              string `json:"item_name"`
	CreateLabelsIfMissing bool   `json:"create_labels_if_missing,omitempty"`
}

type CreateItemAction struct{}

func (a *CreateItemAction) Name() string {
	return "Create Item"
}

func (a *CreateItemAction) Description() string {
	return "Create Item: Automatically generates a new item in your system, such as a task, issue, or project, with customizable fields and attributes."
}

func (a *CreateItemAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateItemAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createItemDocs,
	}
}

func (a *CreateItemAction) Icon() *string {
	return nil
}

func (a *CreateItemAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace_id": shared.GetWorkspaceInput(),
		"board_id":     shared.GetBoardInput("Board ID", "Select Board"),
		"group_id":     shared.GetGroupInput("Group", "Select Group", false),
		"item_name": autoform.NewShortTextField().
			SetDisplayName("Item Name").
			SetDescription("Item Name").
			SetRequired(true).
			Build(),
		"create_labels_if_missing": autoform.NewBooleanField().
			SetDisplayName("Create Labels if Missing").
			SetDescription("Creates status/dropdown labels if they are missing. This requires permission to change the board structure.").
			SetDefaultValue(false).
			SetRequired(false).
			Build(),
	}
}

func (a *CreateItemAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createItemActionProps](ctx.BaseContext)
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

	response, err := shared.MondayClient(ctx.BaseContext, query)
	if err != nil {
		return nil, err
	}

	item, ok := response["data"].(map[string]interface{})["create_item"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract item from response")
	}

	return item, nil
}

func (a *CreateItemAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateItemAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateItemAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateItemAction() sdk.Action {
	return &CreateItemAction{}
}
