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
	"github.com/wakflo/extensions/internal/integrations/square/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type getPaymentsActionProps struct {
	OrganizationID string `json:"organization_id"`
}

type GetPaymentsAction struct{}

func (c *GetPaymentsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c GetPaymentsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c GetPaymentsAction) Name() string {
	return "Get Payments"
}

func (c GetPaymentsAction) Description() string {
	return "Retrieve a list of Payments"
}

func (c GetPaymentsAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &getPaymentDocs,
	}
}

func (c GetPaymentsAction) Icon() *string {
	return nil
}

func (c GetPaymentsAction) SampleData() sdkcore.JSON {
	return nil
}

func (c GetPaymentsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (c GetPaymentsAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c GetPaymentsAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	_, err := integration.InputToTypeSafely[getPaymentsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/v2/payments"

	payments, err := shared.GetSquareClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func NewGetPaymentsAction() integration.Action {
	return &GetPaymentsAction{}
}
