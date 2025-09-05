package actions

import (
	"errors"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listTicketsActionProps struct {
	PerPage *int `json:"per_page,omitempty"`
}

type ListTicketsAction struct{}

// Metadata returns metadata about the action
func (a *ListTicketsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_tickets",
		DisplayName:   "List Tickets",
		Description:   "Retrieve a list of tickets based on filter criteria.",
		Type:          core.ActionTypeAction,
		Documentation: listTicketsDocs,
		Icon:          "mdi:ticket-percent",
		SampleOutput: map[string]any{
			"data": []map[string]any{
				{
					"id":         "123",
					"subject":    "Sample Ticket 1",
					"status":     "2",
					"priority":   "1",
					"created_at": "2023-12-01T12:30:45Z",
				},
				{
					"id":         "124",
					"subject":    "Sample Ticket 2",
					"status":     "3",
					"priority":   "2",
					"created_at": "2023-12-02T10:15:30Z",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListTicketsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_tickets", "List Tickets")

	form.NumberField("per_page", "Per Page").
		Placeholder("Results Per Page").
		HelpText("Number of results per page (max 100)").
		Required(false)

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

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing freshdesk api key")
	}

	endpoint := "/tickets"

	queryParams := "?"

	if input.PerPage != nil {
		queryParams += "per_page=" + strconv.Itoa(*input.PerPage) + "&"
	}

	if queryParams != "?" {
		endpoint += queryParams[:len(queryParams)-1] // Remove trailing & or ?
	}

	domain := authCtx.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	response, err := shared.GetTickets(endpoint, freshdeskDomain, authCtx.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewListTicketsAction() sdk.Action {
	return &ListTicketsAction{}
}
