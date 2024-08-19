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

package square

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getPaymentsOperationProps struct {
	OrganizationID string `json:"organization_id"`
}

type GetPaymentsOperation struct {
	options *sdk.OperationInfo
}

func NewGetPaymentsOperation() *GetPaymentsOperation {
	return &GetPaymentsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Payments",
			Description: "Retrieve a list of Payments",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetPaymentsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Square auth token")
	}

	_ = sdk.InputToType[getPaymentsOperationProps](ctx)

	url := "https://connect.squareup.com/v2/payments"

	payments, err := getSquareClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (c *GetPaymentsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetPaymentsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}