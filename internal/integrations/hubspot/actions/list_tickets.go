package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listTicketsActionProps struct {
	Limit int `json:"limit"`
}

type ListTicketsAction struct{}

func (a *ListTicketsAction) Name() string {
	return "List Tickets"
}

func (a *ListTicketsAction) Description() string {
	return "List available tickets"
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
	return nil
}

func (a *ListTicketsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Limit").
			SetRequired(false).Build(),
	}
}

func (a *ListTicketsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listTicketsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}

	url := fmt.Sprintf("/crm/v3/objects/tickets?limit=%d", input.Limit)

	resp, err := shared.HubspotClient(url, ctx.Auth.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *ListTicketsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListTicketsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"results": []map[string]any{
			{
				"id": "51",
				"properties": map[string]any{
					"firstname": "John",
					"lastname":  "Doe",
					"email":     "john.doe@example.com",
				},
			},
			{
				"id": "52",
				"properties": map[string]any{
					"firstname": "Jane",
					"lastname":  "Smith",
					"email":     "jane.smith@example.com",
				},
			},
		},
		"paging": map[string]any{
			"next": map[string]any{
				"after": "NTI=",
				"link":  "https://api.hubapi.com/crm/v3/objects/contacts?after=NTI=",
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
