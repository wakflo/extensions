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
	"log"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/cin7/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getPaymentActionProps struct {
	PaymentID string `json:"id"`
}

type GetPaymentAction struct{}

// Metadata returns metadata about the action
func (a *GetPaymentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_payment",
		DisplayName:   "Get a payment",
		Description:   "Get payment information",
		Type:          core.ActionTypeAction,
		Documentation: getAPaymentDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetPaymentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_payment", "Get a payment")

	form.TextField("id", "Task ID").
		Required(true).
		HelpText("The task Id of the payment")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetPaymentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetPaymentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getPaymentActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	endpoint := "/ExternalApi/SalePayments"
	accountID := authCtx.Extra["account_id"]
	applicationKey := authCtx.Extra["key"]
	queryParams := map[string]interface{}{
		"taskId": input.PaymentID,
	}

	response, err := shared.FetchData(endpoint, accountID, applicationKey, queryParams)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func NewGetPaymentAction() sdk.Action {
	return &GetPaymentAction{}
}
