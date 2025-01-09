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

type createColumnOperationProps struct {
	ColumnTitle string `json:"column_title,omitempty"`
	BoardID     string `json:"board_id,omitempty"`
	ColumnType  string `json:"column_type,omitempty"`
}

type CreateColumnOperation struct {
	options *sdk.OperationInfo
}

func NewCreateColumnOperation() *CreateColumnOperation {
	return &CreateColumnOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Column",
			Description: "Creates a new column on a board",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace_id": getWorkspaceInput(),
				"board_id":     getBoardInput("Board ID", "Select Board"),
				"column_title": autoform.NewShortTextField().
					SetDisplayName("Column Title").
					SetDescription("Group name").
					SetRequired(true).
					Build(),
				"column_type": autoform.NewSelectField().
					SetDisplayName("Column Type").
					SetOptions(columnType).
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

func (c *CreateColumnOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing monday.com auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input, err := sdk.InputToTypeSafely[createColumnOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	fields["board_id"] = fmt.Sprintf(`"%s"`, input.BoardID)
	fields["title"] = fmt.Sprintf(`"%s"`, input.ColumnTitle)
	fields["column_type"] = input.ColumnType

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	mutation := fmt.Sprintf(`
		mutation {
  			create_column (%s) {
    			id
    			title
    			description
		}
}`, strings.Join(fieldStrings, "\n"))

	response, err := mondayClient(accessToken, mutation)
	if err != nil {
		return nil, err
		// return nil, errors.New("request not successful")
	}

	newColumn, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract column from response")
	}

	return newColumn, nil
}

func (c *CreateColumnOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateColumnOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
