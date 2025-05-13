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

type createTicketActionProps struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	Email       string `json:"email"`
	CCEmails    string `json:"cc_emails"`
}

type CreateTicketAction struct{}

// Metadata returns metadata about the action
func (a *CreateTicketAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_ticket",
		DisplayName:   "Create Ticket",
		Description:   "Create a new support ticket in the Freshdesk system with customizable fields.",
		Type:          core.ActionTypeAction,
		Documentation: createTicketDocs,
		SampleOutput: map[string]any{
			"id":           "123",
			"subject":      "Sample Ticket",
			"description":  "This is a sample ticket created via API",
			"status":       "2",
			"priority":     "1",
			"requester_id": "456",
			"created_at":   "2023-12-01T12:30:45Z",
			"updated_at":   "2023-12-01T12:30:45Z",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateTicketAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_ticket", "Create Ticket")

	// Add subject field
	form.TextField("subject", "Subject").
		Placeholder("Enter a subject").
		Required(true).
		HelpText("The subject of the ticket")

	// Add email field
	form.TextField("email", "Email").
		Placeholder("Enter an email").
		Required(true).
		HelpText("Your email address")

	// Add description field
	form.TextField("description", "Description").
		Placeholder("Enter a description").
		Required(true).
		HelpText("Content of the ticket")

	// Add cc_emails field
	form.TextField("cc_emails", "CC Emails").
		Placeholder("Enter CC emails").
		Required(false).
		HelpText("Email address added in the 'cc' field of the incoming ticket email")

	// Add priority field
	form.SelectField("priority", "Priority").
		Placeholder("Select a priority").
		Required(false).
		HelpText("The priority level of the ticket. The default Value is Low.").
		AddOption("1", "Low").
		AddOption("2", "Normal").
		AddOption("3", "High").
		AddOption("4", "Urgent")

	// Add status field
	form.SelectField("status", "Status").
		Placeholder("Select a status").
		Required(false).
		HelpText("The Status of the ticket. The default Value is Open.").
		AddOption("2", "Open").
		AddOption("3", "Pending").
		AddOption("4", "Resolved").
		AddOption("5", "Closed")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateTicketAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateTicketAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createTicketActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	domain := authCtx.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	ticketData := map[string]interface{}{
		"description": input.Description,
		"subject":     input.Subject,
		"email":       input.Email,
	}

	if input.Status != "" {
		status, err := strconv.Atoi(input.Status)
		if err != nil {
			return nil, err
		}
		ticketData["status"] = status
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}

		ticketData["priority"] = priority
	}

	if input.CCEmails != "" {
		ticketData["cc_emails"] = []string{input.CCEmails}
	}

	response, err := shared.CreateTicket(freshdeskDomain, authCtx.Extra["api-key"], ticketData)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket:  %v", err)
	}

	return response, nil
}

func NewCreateTicketAction() sdk.Action {
	return &CreateTicketAction{}
}
