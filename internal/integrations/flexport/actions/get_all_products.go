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

type getProductsActionProps struct {
	Topic string `json:"topic"`
}

type GetProductsAction struct{}

// Metadata returns metadata about the action
func (a *GetProductsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_products",
		DisplayName:   "Get all products",
		Description:   "Get all products",
		Type:          core.ActionTypeAction,
		Documentation: getProductsDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetProductsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_products", "Get all products")

	form.TextField("topic", "Topic").
		Required(false).
		HelpText("Filter products by topic")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetProductsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetProductsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	_, err := sdk.InputToTypeSafely[getProductsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	reqURL := "/api/2024-07/products"
	resp, err := shared.FlexportRequest(authCtx.Extra["api-key"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewGetProductsAction() sdk.Action {
	return &GetProductsAction{}
}
