package actions

import (
	"fmt"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateTicketActionProps struct {
	TicketID    string `json:"ticket_id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
}

type UpdateTicketAction struct{}

// Metadata returns metadata about the action
func (a *UpdateTicketAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_ticket",
		DisplayName:   "Update Ticket",
		Description:   "Update the properties and fields of an existing Freshdesk ticket.",
		Type:          core.ActionTypeAction,
		Documentation: updateTicketDocs,
		Icon:          "mdi:ticket-edit",
		SampleOutput: map[string]any{
			"id":           "123",
			"subject":      "Updated Ticket Subject",
			"description":  "This ticket has been updated via API",
			"status":       "3",
			"priority":     "2",
			"requester_id": "456",
			"updated_at":   "2023-12-01T15:45:30Z",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateTicketAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_ticket", "Update Ticket")

	form.NumberField("ticket_id", "ticket_id").
		Placeholder("Ticket ID").
		HelpText("The ID of the ticket to update").
		Required(true)

	form.TextField("subject", "subject").
		Placeholder("Subject").
		HelpText("The updated subject of the ticket").
		Required(false)

	form.TextareaField("description", "description").
		Placeholder("Description").
		HelpText("The updated description of the ticket").
		Required(false)

	form.SelectField("priority", "priority").
		Placeholder("Priority").
		HelpText("The priority level of the ticket. The default Value is Low.").
		AddOptions(shared.FreshdeskPriorityType...).
		Required(false)

	form.SelectField("status", "status").
		Placeholder("Status").
		HelpText("The Status of the ticket. The default Value is Open.").
		AddOptions(shared.FreshdeskStatusType...).
		Required(false)

	form.TextField("type", "type").
		Placeholder("Type").
		HelpText("The updated type of the ticket").
		Required(false)

	form.TextField("tags", "tags").
		Placeholder("Tags").
		HelpText("Comma-separated list of updated tags").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateTicketAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateTicketAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTicketActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	domain := authCtx.Extra["domain"]
	// freshdeskDomain := "https://" + domain + ".freshdesk.com"
	freshdeskDomain := shared.BuildFreshdeskURL(domain)

	ticketData := shared.TicketUpdate{}
	ticketID := fmt.Sprintf("%v", input.TicketID)

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

	err = shared.UpdateTicket(freshdeskDomain, authCtx.Extra["api-key"], ticketID, ticketData)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket:  %v", err)
	}

	return core.JSON(map[string]interface{}{
		"Status": "Ticket successfully updated",
	}), nil
}

func NewUpdateTicketAction() sdk.Action {
	return &UpdateTicketAction{}
}
