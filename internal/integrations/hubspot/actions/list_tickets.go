package actions

import (
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listTicketsActionProps struct {
	Limit int `json:"limit"`
}

type ListTicketsAction struct{}

// Metadata returns metadata about the action
func (a *ListTicketsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_tickets",
		DisplayName:   "List Tickets",
		Description:   "List available tickets",
		Type:          core.ActionTypeAction,
		Documentation: listTicketsDocs,
		SampleOutput: map[string]any{
			"results": []map[string]any{
				{
					"id": "101",
					"properties": map[string]any{
						"subject":            "Technical issue with login",
						"content":            "User unable to access the system.",
						"hs_ticket_priority": "HIGH",
						"createdate":         "2023-04-15T09:30:00Z",
					},
				},
				{
					"id": "102",
					"properties": map[string]any{
						"subject":            "Feature request",
						"content":            "Customer would like to request a new reporting feature.",
						"hs_ticket_priority": "MEDIUM",
						"createdate":         "2023-04-16T14:45:00Z",
					},
				},
			},
			"paging": map[string]any{
				"next": map[string]any{
					"after": "MTAy",
					"link":  "https://api.hubapi.com/crm/v3/objects/tickets?after=MTAy",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListTicketsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_tickets", "List Tickets")

	form.NumberField("limit", "Limit").
		Required(false).
		HelpText("Limit")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListTicketsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListTicketsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listTicketsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}

	url := fmt.Sprintf("/crm/v3/objects/tickets?limit=%d", input.Limit)

	resp, err := shared.HubspotClient(url, authCtx.Token.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewListTicketsAction() sdk.Action {
	return &ListTicketsAction{}
}
