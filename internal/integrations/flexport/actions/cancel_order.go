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
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/flexport/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type cancelOrderActionProps struct {
	OrderID string `json:"Id"`
}

type CancelOrderAction struct{}

// Metadata returns metadata about the action
func (a *CancelOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "cancel_order",
		DisplayName:   "Cancel Order",
		Description:   "Cancel Order",
		Type:          core.ActionTypeAction,
		Documentation: cancelOrderDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CancelOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("cancel_order", "Cancel Order")

	form.TextField("Id", "Order ID").
		Required(true).
		HelpText("Order ID")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CancelOrderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CancelOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[cancelOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("/api/2024-07/orders/%s/cancel", input.OrderID)
	resp, err := shared.FlexportRequest(authCtx.Extra["api-key"], reqURL, http.MethodPost, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCancelOrderAction() sdk.Action {
	return &CancelOrderAction{}
}
