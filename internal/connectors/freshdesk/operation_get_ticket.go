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

package freshdesk

import (
	"errors"
	"log"

	"github.com/wakflo/go-sdk/autoform"
	sdk2 "github.com/wakflo/go-sdk/sdk"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getTicketOperationProps struct {
	TicketID string `json:"id"`
}

type GetTicketOperation struct {
	options *sdk.OperationInfo
}

func NewGetTicketOperation() *GetTicketOperation {
	return &GetTicketOperation{
		options: &sdk.OperationInfo{
			Name:        "Get a ticket",
			Description: "Retrieves a specific ticket",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Ticket ID").
					SetDescription(" The id of the ticket").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetTicketOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshdesk auth values")
	}

	input := sdk.InputToType[getTicketOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"
	ticket, err := GetTicket(freshdeskDomain, ctx.Auth.Extra["api-key"], input.TicketID)
	if err != nil {
		log.Fatalf("error fetching data: %v", err)
	}

	return ticket, nil
}

func (c *GetTicketOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetTicketOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
