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

package clickup

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createSpaceOperationProps struct {
	TeamID    string `json:"team-id"`
	SpaceName string `json:"space-name"`
}

type CreateSpaceOperation struct {
	options *sdk.OperationInfo
}

func NewCreateSpaceOperation() *CreateSpaceOperation {
	return &CreateSpaceOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Space",
			Description: "Create a Space for current Team",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"team-id": getWorkSpaceInput("Workspaces", "select workspace", true),
				"space-name": autoform.NewShortTextField().
					SetDisplayName("Space Name").
					SetDescription("The space name").
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

func (c *CreateSpaceOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[createSpaceOperationProps](ctx)

	reqURL := "https://api.clickup.com/api/v2/team/" + input.TeamID + "/space"
	data := []byte(fmt.Sprintf(`{
		"name": "%s",
		"multiple_assignees": true,
		"features": {
			"due_dates": {
				"enabled": true,
				"start_date": false,
				"remap_due_dates": true,
				"remap_closed_due_date": false
			},
			"time_tracking": {
				"enabled": false
			},
			"tags": {
				"enabled": true
			},
			"time_estimates": {
				"enabled": true
			},
			"checklists": {
				"enabled": true
			},
			"custom_fields": {
				"enabled": true
			},
			"remap_dependencies": {
				"enabled": true
			},
			"dependency_warning": {
				"enabled": true
			},
			"portfolios": {
				"enabled": true
			}
		}
	}`, input.SpaceName))
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(res)
	fmt.Println(string(body))

	return map[string]interface{}{
		"Result": "Space created Successfully",
	}, nil
}

func (c *CreateSpaceOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateSpaceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
