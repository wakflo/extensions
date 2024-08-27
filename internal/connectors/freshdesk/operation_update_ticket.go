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
	"fmt"
	"strconv"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateTicketOperationProps struct {
	TicketID    string `json:"ticketId"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	CCEmails    string `json:"cc_emails"`
}

type UpdateTicketOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateTicketOperation() *UpdateTicketOperation {
	return &UpdateTicketOperation{
		options: &sdk.OperationInfo{
			Name:        "Update A Ticket",
			Description: "updates a ticket",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"ticketId": getTicketInput(),
				"subject": autoform.NewShortTextField().
					SetDisplayName(" Subject").
					SetDescription("The subject of the ticket").
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("Content of the ticket").
					Build(),
				"priority": autoform.NewSelectField().
					SetDisplayName("Priority").
					SetDescription("The priority level of the ticket. The default Value is Low.").
					SetOptions(freshdeskPriorityType).
					Build(),
				"status": autoform.NewSelectField().
					SetDisplayName("Status").
					SetDescription("The Status of the ticket. The default Value is Open.").
					SetOptions(freshdeskStatusType).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UpdateTicketOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshdesk auth values")
	}

	input := sdk.InputToType[updateTicketOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	priority, err := strconv.Atoi(input.Priority)
	if err != nil {
		return nil, err
	}

	status, err := strconv.Atoi(input.Status)
	if err != nil {
		return nil, err
	}

	ticketData := TicketUpdate{
		Description: input.Description,
		Subject:     input.Subject,
		Priority:    priority,
		Status:      status,
	}

	err = UpdateTicket(freshdeskDomain, ctx.Auth.Extra["api-key"], input.TicketID, ticketData)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket:  %v", err)
	}

	return sdk.JSON(map[string]interface{}{
		"Status": "Ticket successfully updated",
	}), nil
}

func (c *UpdateTicketOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateTicketOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
