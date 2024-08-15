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
)

type getSpacesOperationProps struct {
	TeamID string `json:"team-id"`
}

type GetSpacesOperation struct {
	options *sdk.OperationInfo
}

func NewGetSpacesOperation() *GetSpacesOperation {
	getAuthTeams := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		authorizedTeams, err := getTeams(ctx.Auth.AccessToken)
		if err != nil {
			return nil, err
		}

		return authorizedTeams, nil
	}
	return &GetSpacesOperation{
		options: &sdk.OperationInfo{
			Name:        "Gets Spaces in a ClickUp workspace",
			Description: "Get Spaces",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"team-id": autoform.NewDynamicField(sdkcore.String).
					SetDisplayName("Team ID").
					SetDescription("Team ID to get the spaces on").
					SetDynamicOptions(&getAuthTeams).
					SetDependsOn([]string{"connection"}).
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

func (c *GetSpacesOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[getSpacesOperationProps](ctx)

	spaces, _ := getSpaces(accessToken, input.TeamID)

	return map[string]interface{}{
		"Spaces": spaces,
	}, nil
}

func (c *GetSpacesOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetSpacesOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
