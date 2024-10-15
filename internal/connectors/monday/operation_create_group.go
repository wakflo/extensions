// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monday

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createGroupOperationProps struct {
	GroupName string `json:"group_name,omitempty"`
	BoardID   string `json:"board_id,omitempty"`
}

type CreateGroupOperation struct {
	options *sdk.OperationInfo
}

func NewCreateGroupOperation() *CreateGroupOperation {
	return &CreateGroupOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Group",
			Description: "Create a group ",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace_id": getWorkspaceInput(),
				"board_id":     getBoardInput("Board ID", "Select Board"),
				"group_name": autoform.NewShortTextField().
					SetDisplayName("Group Name").
					SetDescription("Group name").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateGroupOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing monday.com auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input, err := sdk.InputToTypeSafely[createGroupOperationProps](ctx)
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

	response, err := mondayClient(accessToken, mutation)
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

func (c *CreateGroupOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateGroupOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
