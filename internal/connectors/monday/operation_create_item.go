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

type createItemOperationProps struct {
	WorkspaceID           string `json:"workspace_id,omitempty"`
	BoardID               string `json:"board_id,omitempty"`
	GroupID               string `json:"group_id,omitempty"`
	ItemName              string `json:"item_name"`
	CreateLabelsIfMissing bool   `json:"create_labels_if_missing,omitempty"`
}

type CreateItemOperation struct {
	options *sdk.OperationInfo
}

func NewCreateItemOperation() *CreateItemOperation {
	return &CreateItemOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Item",
			Description: "Creates a new column on a board",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace_id": getWorkspaceInput(),
				"board_id":     getBoardInput("Board ID", "Select Board"),
				"group_id":     getGroupInput("Group", "Select Group", false),
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
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateItemOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing monday.com auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input, err := sdk.InputToTypeSafely[createItemOperationProps](ctx)
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

	response, err := mondayClient(accessToken, query)
	if err != nil {
		return nil, err
	}

	item, ok := response["data"].(map[string]interface{})["create_item"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract item from response")
	}

	return item, nil
}

func (c *CreateItemOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateItemOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
