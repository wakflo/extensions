package actions

import (
	"fmt"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *CreateTicketAction) Name() string {
	return "Create Ticket"
}

func (a *CreateTicketAction) Description() string {
	return "Create a new support ticket in the Freshdesk system with customizable fields."
}

func (a *CreateTicketAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateTicketAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createTicketDocs,
	}
}

func (a *CreateTicketAction) Icon() *string {
	icon := "mdi:ticket-outline"
	return &icon
}

func (a *CreateTicketAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"subject": autoform.NewShortTextField().
			SetDisplayName(" Subject").
			SetDescription("The subject of the ticket").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName(" Email").
			SetDescription("Your email address").
			SetRequired(true).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("Content of the ticket").
			SetRequired(true).
			Build(),
		"cc_emails": autoform.NewLongTextField().
			SetDisplayName("CC Emails").
			SetDescription(" Email address added in the 'cc' field of the incoming ticket email").
			Build(),
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
	}
}

func (a *CreateTicketAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTicketActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	domain := ctx.Auth.Extra["domain"]
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

	response, err := shared.CreateTicket(freshdeskDomain, ctx.Auth.Extra["api-key"], ticketData)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket:  %v", err)
	}

	return response, nil
}

func (a *CreateTicketAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateTicketAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":           "123",
		"subject":      "Sample Ticket",
		"description":  "This is a sample ticket created via API",
		"status":       "2",
		"priority":     "1",
		"requester_id": "456",
		"created_at":   "2023-12-01T12:30:45Z",
		"updated_at":   "2023-12-01T12:30:45Z",
	}
}

func (a *CreateTicketAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateTicketAction() sdk.Action {
	return &CreateTicketAction{}
}
