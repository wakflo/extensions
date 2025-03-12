package actions

import (
	"errors"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listTicketsActionProps struct {
	Filter  *string `json:"filter,omitempty"`
	Page    *int    `json:"page,omitempty"`
	PerPage *int    `json:"per_page,omitempty"`
}

type ListTicketsAction struct{}

func (a *ListTicketsAction) Name() string {
	return "List Tickets"
}

func (a *ListTicketsAction) Description() string {
	return "Retrieve a list of tickets based on filter criteria."
}

func (a *ListTicketsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListTicketsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listTicketsDocs,
	}
}

func (a *ListTicketsAction) Icon() *string {
	icon := "mdi:ticket-percent"
	return &icon
}

func (a *ListTicketsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"filter": autoform.NewSelectField().
			SetDisplayName("Filter").
			SetDescription("Filter tickets by predefined filters").
			SetRequired(false).
			SetOptions([]*sdkcore.AutoFormSchema{
				{Title: "All Tickets", Const: "all_tickets"},
				{Title: "Open Tickets", Const: "open"},
				{Title: "Pending Tickets", Const: "pending"},
				{Title: "Resolved Tickets", Const: "resolved"},
				{Title: "Closed Tickets", Const: "closed"},
				{Title: "New and My Open Tickets", Const: "new_and_my_open"},
				{Title: "Watching Tickets", Const: "watching"},
				{Title: "Deleted Tickets", Const: "deleted"},
			}).Build(),
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("Page number for pagination").
			SetRequired(false).
			Build(),
		"per_page": autoform.NewNumberField().
			SetDisplayName("Results Per Page").
			SetDescription("Number of results per page (max 100)").
			SetRequired(false).
			Build(),
	}
}

func (a *ListTicketsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listTicketsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing freshdesk api key")
	}

	endpoint := "/tickets"

	queryParams := "?"
	if input.Filter != nil && *input.Filter != "" {
		queryParams += "filter=" + *input.Filter + "&"
	}

	if input.Page != nil {
		queryParams += "page=" + strconv.Itoa(*input.Page) + "&"
	}

	if input.PerPage != nil {
		queryParams += "per_page=" + strconv.Itoa(*input.PerPage) + "&"
	}

	if queryParams != "?" {
		endpoint += queryParams[:len(queryParams)-1] // Remove trailing & or ?
	}

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	response, err := shared.GetTickets(endpoint, freshdeskDomain, ctx.Auth.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *ListTicketsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListTicketsAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *ListTicketsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListTicketsAction() sdk.Action {
	return &ListTicketsAction{}
}
