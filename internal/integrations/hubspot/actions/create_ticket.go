package actions

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createTicketActionProps struct {
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Priority string `json:"hs_ticket_priority"`
}

type CreateTicketAction struct{}

func (a *CreateTicketAction) Name() string {
	return "Create Ticket"
}

func (a *CreateTicketAction) Description() string {
	return "Create a new ticket in HubSpot."
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
	return nil
}

func (a *CreateTicketAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"subject": autoform.NewShortTextField().
			SetDisplayName("Ticket Subject").
			SetDescription("subject of the ticket to create").
			SetRequired(true).Build(),
		"content": autoform.NewShortTextField().
			SetDisplayName("Ticket Description").
			SetDescription("ticket description").
			SetRequired(false).Build(),
		"hs_ticket_priority": autoform.NewSelectField().
			SetDisplayName("Priority").
			SetDescription("The priority of the ticket").
			SetOptions(shared.HubspotPriority).
			SetRequired(false).
			Build(),
	}
}

func (a *CreateTicketAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTicketActionProps](ctx.BaseContext)
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

	resp, err := shared.HubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodPost, newTicket)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *CreateTicketAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateTicketAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateTicketAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateTicketAction() sdk.Action {
	return &CreateTicketAction{}
}
