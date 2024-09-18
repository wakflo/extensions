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

package aftership

import (
	"errors"

	"github.com/aftership/tracking-sdk-go/v5"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type GetUserCouriersOperation struct {
	options *sdk.OperationInfo
}

func NewGetUserCouriersOperation() *GetUserCouriersOperation {
	return &GetUserCouriersOperation{
		options: &sdk.OperationInfo{
			Name:        "Get User Activated Couriers",
			Description: "get couriers activated by user",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetUserCouriersOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}
	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		return nil, err
	}
	result, err := afterShipSdk.Courier.GetUserCouriers().Execute()
	if err != nil {
		return nil, err
	}

	return result.Couriers, nil
}

func (c *GetUserCouriersOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetUserCouriersOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
