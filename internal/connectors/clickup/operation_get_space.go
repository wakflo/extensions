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

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type getSpaceOperationProps struct {
	SpaceID string `json:"space-id"`
}

type GetSpaceOperation struct {
	options *sdk.OperationInfo
}

func NewGetSpaceOperation() *GetSpaceOperation {
	return &GetSpaceOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Space",
			Description: "Gets a space in a ClickUp",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace-id": getWorkSpaceInput("Workspaces", "Get workspaces", true),
				"space-id":     getSpacesInput("Spaces", "get a space", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetSpaceOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[getSpaceOperationProps](ctx)

	space, err := getSpace(accessToken, input.SpaceID)
	if err != nil {
		return nil, err
	}

	return space, nil
}

func (c *GetSpaceOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetSpaceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
