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

package hubspot

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type createNewTicketProps struct {
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Priority string `json:"hs_ticket_priority"`
}

type CreateTicketOperation struct {
	options *sdk.OperationInfo
}

func NewCreateTicketOperation() *CreateTicketOperation {
	return &CreateTicketOperation{
		options: &sdk.OperationInfo{
			Name:        "Create new Ticket",
			Description: "Create new contact",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"subject": autoform.NewShortTextField().
					SetDisplayName("Ticket Subject").
					SetDescription("subject of the ticket to create").
					SetRequired(true).Build(),
				"content": autoform.NewShortTextField().
					SetDisplayName("Ticket Description").
					SetDescription("ticket description").
					SetRequired(false).Build(),
				"hs_ticket_priority": autoform.NewSelectField().
					SetDisplayName("Format").
					SetDescription("The format of the email to read").
					SetOptions(hubspotPriority).
					SetRequired(false).
					Build(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *CreateTicketOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[createNewTicketProps](ctx)

	ticket := contactRequest{
		Properties: map[string]interface{}{
			"subject":            input.Subject,
			"content":            input.Content,
			"hs_ticket_priority": input.Priority,
			"hs_pipeline":        "0",
			"hs_pipeline_stage":  "1",
		},
	}

	newTicket, err := json.Marshal(ticket)
	if err != nil {
		return nil, err
	}

	reqURL := "/crm/v3/objects/tickets"

	resp, err := hubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodPost, newTicket)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CreateTicketOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateTicketOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
