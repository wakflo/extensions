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
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type searchTaskOperationProps struct {
	TeamID        string `json:"team-id"`
	Page          int    `json:"page"`
	Reverse       bool   `json:"reverse"`
	IncludeClosed bool   `json:"include-closed"`
	OrderBy       string `json:"order-by"`
}

type SearchTaskOperation struct {
	options *sdk.OperationInfo
}

func NewSearchTaskOperation() *SearchTaskOperation {
	return &SearchTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Search Team tasks",
			Description: "Retrieves the tasks that meet specific criteria from a Workspace.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"team-id": getWorkSpaceInput("Workspaces", "select workspace", true),
				"page": autoform.NewNumberField().
					SetDisplayName("Page").
					SetDescription("Page to fetch (starts at 0).").
					Build(),
				"reverse": autoform.NewBooleanField().
					SetDisplayName("Reverse").
					SetDescription("Tasks are displayed in reverse order.").
					SetDefaultValue(false).
					Build(),
				"include-closed": autoform.NewBooleanField().
					SetDisplayName("Include Closed").
					SetDescription("Include or exclude closed tasks. By default, they are excluded.").
					SetDefaultValue(false).
					Build(),
				"order-by": autoform.NewSelectField().
					SetDisplayName("Order by").
					SetDescription("Order by a particular field. By default, tasks are ordered by created.").
					SetOptions(clickupOrderbyType).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *SearchTaskOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[searchTaskOperationProps](ctx)

	reqURL := "/v2/team/" + input.TeamID + "/task"
	task, _ := searchTask(accessToken, reqURL, input.Page, input.OrderBy, input.Reverse, input.IncludeClosed)

	return task, nil
}

func (c *SearchTaskOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *SearchTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
