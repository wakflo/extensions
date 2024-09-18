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

type GetCouriersOperation struct {
	options *sdk.OperationInfo
}

func NewGetCouriersOperation() *GetCouriersOperation {
	return &GetCouriersOperation{
		options: &sdk.OperationInfo{
			Name:        "Get All Couriers",
			Description: "get all couriers",
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

func (c *GetCouriersOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}
	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		return nil, err
	}
	result, err := afterShipSdk.Courier.GetAllCouriers().Execute()
	if err != nil {
		return nil, err
	}
	return result.Couriers, nil
}

func (c *GetCouriersOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetCouriersOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
