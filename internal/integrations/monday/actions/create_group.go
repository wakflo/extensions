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

type createGroupActionProps struct {
	GroupName string `json:"group_name,omitempty"`
	BoardID   string `json:"board_id,omitempty"`
}

type CreateGroupAction struct{}

// func (a *CreateGroupAction) Name() string {
// 	return "Create Group"
// }

func (a *CreateGroupAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_group",
		DisplayName:   "Create Group",
		Description:   "Create Group: Creates a new group in your organization's directory, allowing you to categorize and manage users with similar roles or responsibilities. This integration action enables you to automate the process of creating groups, streamlining your workflow and improving collaboration among team members.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createGroupDocs,
		SampleOutput: map[string]any{
			"data": map[string]interface{}{
				"id":   "123456789",
				"name": "New Group",
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateGroupAction) Description() string {
	return "Create Group: Creates a new group in your organization's directory, allowing you to categorize and manage users with similar roles or responsibilities. This integration action enables you to automate the process of creating groups, streamlining your workflow and improving collaboration among team members."
}

func (a *CreateGroupAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_group", "Create Group")

	shared.GetWorkspaceProp(form)
	shared.GetBoardProp("board_id", "Board ID", "Select Board", form)
	form.TextField("group_name", "Group Name").
		Placeholder("New Group").
		Required(true).
		HelpText("Group name.")

	schema := form.Build()

	return schema
}

func (a *CreateGroupAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createGroupActionProps](ctx)
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

	response, err := shared.MondayClient(ctx, mutation)
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

func (a *CreateGroupAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateGroupAction() sdk.Action {
	return &CreateGroupAction{}
}
