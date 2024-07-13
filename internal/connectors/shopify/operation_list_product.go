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

package shopify

import (
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type ListProductsOperation struct {
	options *sdk.OperationInfo
}

func NewListProductsOperation() *ListProductsOperation {
	return &ListProductsOperation{
		options: &sdk.OperationInfo{
			Name:        "Create New File",
			Description: "operation creates new google drive file",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId": autoform.NewLongTextField().
					SetDisplayName("Project ID").
					SetDescription("project id").
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

func (c *ListProductsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	auth, err := ctx.Auth.GetCustomAuth()
	if err != nil {
		return nil, err
	}

	// fake implementation to prove test
	return auth, nil
}

func (c *ListProductsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListProductsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}