package actions

import (
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getTicketActionProps struct {
	TicketID string `json:"id"`
}

type GetTicketAction struct{}

func (a *GetTicketAction) Name() string {
	return "Get Ticket"
}

func (a *GetTicketAction) Description() string {
	return "Retrieve detailed information about a specific ticket by its ID."
}

func (a *GetTicketAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetTicketAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getTicketDocs,
	}
}

func (a *GetTicketAction) Icon() *string {
	icon := "mdi:ticket-account"
	return &icon
}

func (a *GetTicketAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Ticket ID").
			SetDescription(" The id of the ticket").
			SetRequired(true).
			Build(),
	}
}

func (a *GetTicketAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTicketActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"
	ticket, err := shared.GetTicket(freshdeskDomain, ctx.Auth.Extra["api-key"], input.TicketID)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (a *GetTicketAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetTicketAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *GetTicketAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetTicketAction() sdk.Action {
	return &GetTicketAction{}
}
