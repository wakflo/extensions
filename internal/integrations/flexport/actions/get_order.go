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
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/flexport/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

// Documentation for the Get Order by ID action

type getOrderActionProps struct {
	OrderID string `json:"Id"`
}

type GetOrderAction struct{}

// Metadata returns metadata about the action
func (a *GetOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_order",
		DisplayName:   "Get order by ID",
		Description:   "Get order by ID",
		Type:          core.ActionTypeAction,
		Documentation: getOrderDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_order", "Get order by ID")

	form.TextField("Id", "Order ID").
		Required(true).
		HelpText("Order ID from flexport")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetOrderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	reqURL := "/api/2024-07/orders/" + input.OrderID
	resp, err := shared.FlexportRequest(authCtx.Extra["api-key"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewGetOrderAction() sdk.Action {
	return &GetOrderAction{}
}
