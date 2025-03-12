package actions

import (
	"fmt"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateTicketActionProps struct {
	TicketID    string `json:"ticket_id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	Tags        string `json:"tags"`
	Type        string `json:"type"`
}

type UpdateTicketAction struct{}

func (a *UpdateTicketAction) Name() string {
	return "Update Ticket"
}

func (a *UpdateTicketAction) Description() string {
	return "Update the properties and fields of an existing Freshdesk ticket."
}

func (a *UpdateTicketAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateTicketAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateTicketDocs,
	}
}

func (a *UpdateTicketAction) Icon() *string {
	icon := "mdi:ticket-edit"
	return &icon
}

func (a *UpdateTicketAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"ticket_id": autoform.NewNumberField().
			SetDisplayName("Ticket ID").
			SetDescription("The ID of the ticket to update").
			SetRequired(true).Build(),
		"subject": autoform.NewShortTextField().
			SetDisplayName("Subject").
			SetDescription("The updated subject of the ticket").
			SetRequired(false).Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("The updated description of the ticket").
			SetRequired(false).Build(),
		"priority": autoform.NewSelectField().
			SetDisplayName("Priority").
			SetDescription("The priority level of the ticket. The default Value is Low.").
			SetOptions(shared.FreshdeskPriorityType).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetDescription("The Status of the ticket. The default Value is Open.").
			SetOptions(shared.FreshdeskStatusType).
			Build(),
		"type": autoform.NewShortTextField().
			SetDisplayName("Type").
			SetDescription("The updated type of the ticket").
			SetRequired(false).Build(),
		"tags": autoform.NewShortTextField().
			SetDisplayName("Tags").
			SetDescription("Comma-separated list of updated tags").
			SetRequired(false).Build(),
	}
}

func (a *UpdateTicketAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTicketActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	ticketData := shared.TicketUpdate{}

	if input.Subject != "" {
		ticketData.Subject = input.Subject
	}
	if input.Description != "" {
		ticketData.Description = input.Description
	}
	if input.Status != "" {
		status, err := strconv.Atoi(input.Status)
		if err != nil {
			return nil, err
		}
		ticketData.Status = status
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}
		ticketData.Priority = priority
	}

	if input.Tags != "" {
		ticketData.Tags = input.Tags
	}

	if input.Type != "" {
		ticketData.Type = input.Type
	}

	err = shared.UpdateTicket(freshdeskDomain, ctx.Auth.Extra["api-key"], input.TicketID, ticketData)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket:  %v", err)
	}

	return sdk.JSON(map[string]interface{}{
		"Status": "Ticket successfully updated",
	}), nil
}

func (a *UpdateTicketAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateTicketAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":           "123",
		"subject":      "Updated Ticket Subject",
		"description":  "This ticket has been updated via API",
		"status":       "3",
		"priority":     "2",
		"requester_id": "456",
		"updated_at":   "2023-12-01T15:45:30Z",
	}
}

func (a *UpdateTicketAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateTicketAction() sdk.Action {
	return &UpdateTicketAction{}
}
