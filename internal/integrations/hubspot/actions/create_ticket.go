package actions

import (
	"encoding/json"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createTicketActionProps struct {
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Priority string `json:"hs_ticket_priority"`
}

type CreateTicketAction struct{}

// Metadata returns metadata about the action
func (a *CreateTicketAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_ticket",
		DisplayName:   "Create Ticket",
		Description:   "Create a new ticket in HubSpot.",
		Type:          core.ActionTypeAction,
		Documentation: createTicketDocs,
		SampleOutput: map[string]any{
			"id": "12345",
			"properties": map[string]any{
				"subject":            "Technical issue with product",
				"content":            "Customer reported an issue with logging in to the application.",
				"hs_ticket_priority": "HIGH",
				"hs_pipeline":        "0",
				"hs_pipeline_stage":  "1",
				"createdate":         "2023-05-01T12:00:00.000Z",
			},
			"createdAt": "2023-05-01T12:00:00.000Z",
			"updatedAt": "2023-05-01T12:00:00.000Z",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateTicketAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_ticket", "Create Ticket")

	form.TextField("subject", "Ticket Subject").
		Required(true).
		HelpText("subject of the ticket to create")

	form.TextareaField("content", "Ticket Description").
		Required(false).
		HelpText("ticket description")

	form.SelectField("hs_ticket_priority", "Priority").
		Required(false).
		AddOption("LOW", "Low").
		AddOption("MEDIUM", "Medium").
		AddOption("HIGH", "High").
		HelpText("The priority of the ticket")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateTicketAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateTicketAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTicketActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	ticket := shared.ContactRequest{
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

	resp, err := shared.HubspotClient(reqURL, authCtx.Token.AccessToken, http.MethodPost, newTicket)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCreateTicketAction() sdk.Action {
	return &CreateTicketAction{}
}
