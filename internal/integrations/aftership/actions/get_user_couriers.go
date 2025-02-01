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
	"github.com/wakflo/go-sdk/integration"
)

type GetUserCouriersAction struct{}

func (c GetUserCouriersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c GetUserCouriersAction) Name() string {
	return "Get User Activated Couriers"
}

func (c GetUserCouriersAction) Description() string {
	return "get couriers activated by user"
}

func (c GetUserCouriersAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &getUserCouriersDocs,
	}
}

func (c GetUserCouriersAction) Icon() *string {
	return nil
}

func (c GetUserCouriersAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c GetUserCouriersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (c GetUserCouriersAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c GetUserCouriersAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
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

func NewGetUserCouriersAction() integration.Action {
	return &GetUserCouriersAction{}
}
