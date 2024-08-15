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

type getListOperationProps struct {
	ListID string `json:"list-id"`
}

type GetListOperation struct {
	options *sdk.OperationInfo
}

func NewGetListOperation() *GetListOperation {
	return &GetListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get a list",
			Description: "The ID of the List",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"list-id": autoform.NewShortTextField().
					SetDisplayName("List ID").
					SetDescription("List ID").
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

func (c *GetListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[getListOperationProps](ctx)

	list, _ := getList(accessToken, input.ListID)

	return list, nil
}

func (c *GetListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
