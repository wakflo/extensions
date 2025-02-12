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

type createGroupActionProps struct {
	GroupName string `json:"group_name,omitempty"`
	BoardID   string `json:"board_id,omitempty"`
}

type CreateGroupAction struct{}

func (a *CreateGroupAction) Name() string {
	return "Create Group"
}

func (a *CreateGroupAction) Description() string {
	return "Create Group: Creates a new group in your organization's directory, allowing you to categorize and manage users with similar roles or responsibilities. This integration action enables you to automate the process of creating groups, streamlining your workflow and improving collaboration among team members."
}

func (a *CreateGroupAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateGroupAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createGroupDocs,
	}
}

func (a *CreateGroupAction) Icon() *string {
	return nil
}

func (a *CreateGroupAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace_id": shared.GetWorkspaceInput(),
		"board_id":     shared.GetBoardInput("Board ID", "Select Board"),
		"group_name": autoform.NewShortTextField().
			SetDisplayName("Group Name").
			SetDescription("Group name").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateGroupAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createGroupActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	fields["board_id"] = fmt.Sprintf(`"%s"`, input.BoardID)
	fields["group_name"] = fmt.Sprintf(`"%s"`, input.GroupName)

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	mutation := fmt.Sprintf(`
		mutation {
  			create_group (%s) {
    			id
		}
}`, strings.Join(fieldStrings, "\n"))

	response, err := shared.MondayClient(ctx.BaseContext, mutation)
	if err != nil {
		return nil, err
		// return nil, errors.New("request not successful")
	}

	group, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract item from response")
	}

	return group, nil
}

func (a *CreateGroupAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateGroupAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateGroupAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateGroupAction() sdk.Action {
	return &CreateGroupAction{}
}
