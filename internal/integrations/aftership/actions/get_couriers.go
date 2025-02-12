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

package actions

import (
	"errors"

	"github.com/aftership/tracking-sdk-go/v5"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type GetCouriersAction struct{}

func (c *GetCouriersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c GetCouriersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c GetCouriersAction) Name() string {
	return "Get All Couriers"
}

func (c GetCouriersAction) Description() string {
	return "get all couriers"
}

func (c GetCouriersAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getCouriersDocs,
	}
}

func (c GetCouriersAction) Icon() *string {
	return nil
}

func (c GetCouriersAction) SampleData() sdkcore.JSON {
	return nil
}

func (c GetCouriersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (c GetCouriersAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c GetCouriersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
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

func NewGetCouriersAction() sdk.Action {
	return &GetCouriersAction{}
}
