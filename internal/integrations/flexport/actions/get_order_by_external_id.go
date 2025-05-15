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

type getOrderByExternalIDActionProps struct {
	ExternalOrderID string `json:"Id"`
}

type GetOrderByExternalIDAction struct{}

// Metadata returns metadata about the action
func (a *GetOrderByExternalIDAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_order_by_external_id",
		DisplayName:   "Get order by External ID",
		Description:   "Get order using the external order ID given during order creation.",
		Type:          core.ActionTypeAction,
		Documentation: getOrderByExternalIDDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetOrderByExternalIDAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_order_by_external_id", "Get order by External ID")

	form.TextField("Id", "External Order ID").
		Required(true).
		HelpText("External Order ID from store(woocommerce, shopify etc.)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetOrderByExternalIDAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetOrderByExternalIDAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getOrderByExternalIDActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	reqURL := "/api/2024-07/orders/external_id/" + input.ExternalOrderID
	resp, err := shared.FlexportRequest(authCtx.Extra["api-key"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewGetOrderByExternalIDAction() sdk.Action {
	return &GetOrderByExternalIDAction{}
}
