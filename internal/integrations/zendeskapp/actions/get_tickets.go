package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/zendeskapp/shared"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getTicketsActionProps struct {
}

type GetTicketsAction struct{}

func (a *GetTicketsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Get Tickets",
		Description:   "Get Tickets: Retrieves all support tickets from your Zendesk account, allowing you to manage and track customer inquiries and issues.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getTicketsDocs,
		SampleOutput: map[string]any{
			"tickets": []map[string]any{
				{
					"id":      123456,
					"subject": "Help with account setup",
					"status":  "open",
				},
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetTicketsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_tickets", "Get Tickets")

	schema := form.Build()
	return schema
}

func (a *GetTicketsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	// Extract authentication details
	authData := ctx.Auth()
	if authData == nil {
		return nil, errors.New("missing authentication data")
	}

	// Check if required auth details are present
	email, ok := authData.Extra["email"]
	if !ok || email == "" {
		return nil, errors.New("missing zendesk email")
	}

	apiToken, ok := authData.Extra["api-token"]
	if !ok || apiToken == "" {
		return nil, errors.New("missing zendesk api-token")
	}

	subdomain, ok := authData.Extra["subdomain"]
	if !ok || subdomain == "" {
		return nil, errors.New("missing zendesk subdomain")
	}

	// Construct the URL
	fullURL := fmt.Sprintf("https://%s.zendesk.com/api/v2/tickets.json", subdomain)

	// Make the request
	response, err := shared.ZendeskRequest(http.MethodGet, fullURL, email, apiToken, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching tickets: %v", err)
	}

	// Extract tickets from the response
	tickets, ok := response["tickets"]
	if !ok {
		return nil, errors.New("failed to extract tickets from response")
	}

	return tickets, nil
}

func (a *GetTicketsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetTicketsAction() sdk.Action {
	return &GetTicketsAction{}
}
