package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getTicketActionProps struct {
	TicketID string `json:"ticketId"`
}

type GetTicketAction struct{}

// Metadata returns metadata about the action
func (a *GetTicketAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_ticket",
		DisplayName:   "Get Ticket",
		Description:   "Retrieve detailed information about a specific ticket by its ID.",
		Type:          core.ActionTypeAction,
		Documentation: getTicketDocs,
		Icon:          "mdi:ticket-account",
		SampleOutput: map[string]any{
			"id":           "123",
			"subject":      "Sample Ticket",
			"description":  "This is a sample ticket details",
			"status":       "2",
			"priority":     "1",
			"requester_id": "456",
			"responder_id": "789",
			"created_at":   "2023-12-01T12:30:45Z",
			"updated_at":   "2023-12-01T14:20:15Z",
			"due_by":       "2023-12-03T12:30:45Z",
			"fr_due_by":    "2023-12-02T12:30:45Z",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetTicketAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_ticket", "Get Ticket")

	shared.RegisterTicketProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetTicketAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetTicketAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTicketActionProps](ctx)
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

	ticket, err := shared.GetTicket(freshdeskDomain, authCtx.Extra["api-key"], input.TicketID)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func NewGetTicketAction() sdk.Action {
	return &GetTicketAction{}
}
