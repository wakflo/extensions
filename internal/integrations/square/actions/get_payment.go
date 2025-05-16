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
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/square/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getPaymentsActionProps struct {
	OrganizationID string `json:"organization_id"`
}

type GetPaymentsAction struct{}

func (c *GetPaymentsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_payments",
		DisplayName:   "Get Payments",
		Description:   "Retrieve a list of Payments",
		Type:          core.ActionTypeAction,
		Documentation: getPaymentDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

func (c *GetPaymentsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_payments", "Get Payments")

	schema := form.Build()
	return schema
}

func (c *GetPaymentsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	_, err := sdk.InputToTypeSafely[getPaymentsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/v2/payments"

	payments, err := shared.GetSquareClient(authCtx.Token.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (c *GetPaymentsAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetPaymentsAction() sdk.Action {
	return &GetPaymentsAction{}
}
