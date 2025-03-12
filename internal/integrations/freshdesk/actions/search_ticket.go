package actions

import (
	"errors"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/freshdesk/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type searchTicketsActionProps struct {
	Query    *string `json:"query,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PerPage  *int    `json:"per_page,omitempty"`
	FilterBy *string `json:"filter_by,omitempty"`
}

type SearchTicketsAction struct{}

func (a *SearchTicketsAction) Name() string {
	return "Search Tickets"
}

func (a *SearchTicketsAction) Description() string {
	return "Search for tickets using various search criteria in Freshdesk."
}

func (a *SearchTicketsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SearchTicketsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &searchTicketDocs,
	}
}

func (a *SearchTicketsAction) Icon() *string {
	icon := "mdi:ticket-search"
	return &icon
}

func (a *SearchTicketsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"query": autoform.NewLongTextField().
			SetDisplayName("Search Query").
			SetDescription("Search query string (e.g., 'status:open priority:high')").
			SetRequired(true).Build(),

		"filter_by": autoform.NewSelectField().
			SetDisplayName("Filter By").
			SetDescription("Optional filter to narrow down search results").
			SetRequired(false).
			SetOptions([]*sdkcore.AutoFormSchema{
				{Title: "All Tickets", Const: "all_tickets"},
				{Title: "Open Tickets", Const: "open"},
				{Title: "Pending Tickets", Const: "pending"},
				{Title: "Resolved Tickets", Const: "resolved"},
				{Title: "Closed Tickets", Const: "closed"},
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

func (a *SearchTicketsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchTicketsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Validate required inputs
	if input.Query == nil || *input.Query == "" {
		return nil, errors.New("search query is required")
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing Freshdesk API key")
	}

	endpoint := "/search/tickets"

	queryParams := "?"
	queryParams += "query=" + *input.Query + "&"

	if input.FilterBy != nil && *input.FilterBy != "" {
		queryParams += "filter=" + *input.FilterBy + "&"
	}

	if input.Page != nil {
		queryParams += "page=" + strconv.Itoa(*input.Page) + "&"
	}

	if input.PerPage != nil {
		queryParams += "per_page=" + strconv.Itoa(*input.PerPage) + "&"
	}

	// Remove trailing &
	if queryParams[len(queryParams)-1] == '&' {
		queryParams = queryParams[:len(queryParams)-1]
	}

	// Construct full endpoint
	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"

	response, err := shared.GetTickets(endpoint+queryParams, freshdeskDomain, ctx.Auth.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *SearchTicketsAction) Auth() *sdk.Auth {
	return nil
}

func (a *SearchTicketsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"data": []map[string]any{
			{
				"id":         "123",
				"subject":    "Search Result Ticket 1",
				"status":     "2",
				"priority":   "1",
				"created_at": "2023-12-01T12:30:45Z",
			},
			{
				"id":         "124",
				"subject":    "Search Result Ticket 2",
				"status":     "3",
				"priority":   "2",
				"created_at": "2023-12-02T10:15:30Z",
			},
		},
	}
}

func (a *SearchTicketsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSearchTicketsAction() sdk.Action {
	return &SearchTicketsAction{}
}
