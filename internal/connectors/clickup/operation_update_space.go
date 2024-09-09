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

type updateSpaceOperationProps struct {
	SpaceID           string `json:"space-id"`
	SpaceName         string `json:"space-name"`
	MultipleAssignees bool   `json:"multiple-assignees"`
	Tags              bool   `json:"tags"`
	CustomFields      bool   `json:"custom-fields"`
}

type UpdateSpaceOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateSpaceOperation() *UpdateSpaceOperation {
	return &UpdateSpaceOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Space",
			Description: "Update a Space for current Team",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace-id": getWorkSpaceInput("Workspaces", "select a workspace", true),
				"space-id":     getSpacesInput("Spaces", "select a space", true),
				"space-name": autoform.NewShortTextField().
					SetDisplayName("Update space Name").
					SetDescription("The space name to update").
					SetRequired(false).
					Build(),
				"multiple-assignees": autoform.NewBooleanField().
					SetDisplayName("Multiple Assignees").
					SetDescription("Enable multiple assignees").
					SetRequired(false).
					SetDefaultValue(true).
					Build(),
				"custom-fields": autoform.NewBooleanField().
					SetDisplayName("Custom fields").
					SetDescription("Enable custom fields").
					SetRequired(false).
					SetDefaultValue(true).
					Build(),
				"tags": autoform.NewBooleanField().
					SetDisplayName("Tags").
					SetDescription("Enable tags").
					SetRequired(false).
					SetDefaultValue(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UpdateSpaceOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[updateSpaceOperationProps](ctx)

	reqURL := "https://api.clickup.com/api/v2/space/" + input.SpaceID
	data := []byte(fmt.Sprintf(`{
		"name": "%s",
		"multiple_assignees": %t,
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
				"enabled": %t
			},
			"time_estimates": {
				"enabled": true
			},
			"checklists": {
				"enabled": true
			},
			"custom_fields": {
				"enabled": %t
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
	}`, input.SpaceName, input.MultipleAssignees, input.Tags, input.CustomFields))
	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(data))
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
		"Result": "Space updated Successfully",
	}, nil
}

func (c *UpdateSpaceOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateSpaceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
