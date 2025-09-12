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

type searchTicketsActionProps struct {
	Query    *string `json:"query,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PerPage  *int    `json:"per_page,omitempty"`
	FilterBy *string `json:"filter_by,omitempty"`
}

type SearchTicketsAction struct{}

// Metadata returns metadata about the action
func (a *SearchTicketsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "search_tickets",
		DisplayName:   "Search Tickets",
		Description:   "Search for tickets using various search criteria in Freshdesk.",
		Type:          core.ActionTypeAction,
		Documentation: searchTicketDocs,
		Icon:          "mdi:ticket-search",
		SampleOutput: map[string]any{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *SearchTicketsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("search_tickets", "Search Tickets")

	form.TextareaField("query", "query").
		Placeholder("Search Query").
		HelpText("Search query string (e.g., 'status:open priority:high')").
		Required(true)

	form.SelectField("filter_by", "filter_by").
		Placeholder("Filter By").
		HelpText("Optional filter to narrow down search results").
		AddOptions([]*smartform.Option{
			{Label: "All Tickets", Value: "all_tickets"},
			{Label: "Open Tickets", Value: "open"},
			{Label: "Pending Tickets", Value: "pending"},
			{Label: "Resolved Tickets", Value: "resolved"},
			{Label: "Closed Tickets", Value: "closed"},
		}...).
		Required(false)

	form.NumberField("page", "page").
		Placeholder("Page").
		HelpText("Page number for pagination").
		Required(false)

	form.NumberField("per_page", "per_page").
		Placeholder("Results Per Page").
		HelpText("Number of results per page (max 100)").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *SearchTicketsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *SearchTicketsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchTicketsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Validate required inputs
	if input.Query == nil || *input.Query == "" {
		return nil, errors.New("search query is required")
	}

	if authCtx.Extra["api-key"] == "" {
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
	domain := authCtx.Extra["domain"]
	// freshdeskDomain := "https://" + domain + ".freshdesk.com"
	freshdeskDomain := shared.BuildFreshdeskURL(domain)

	response, err := shared.GetTickets(endpoint+queryParams, freshdeskDomain, authCtx.Extra["api-key"])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewSearchTicketsAction() sdk.Action {
	return &SearchTicketsAction{}
}
