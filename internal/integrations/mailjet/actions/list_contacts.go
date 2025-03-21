package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listContactsActionProps struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type ListContactsAction struct{}

func (a *ListContactsAction) Name() string {
	return "List Contacts"
}

func (a *ListContactsAction) Description() string {
	return "Retrieve contacts from your Mailjet account with optional filtering."
}

func (a *ListContactsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListContactsAction) Icon() *string {
	icon := "mdi:account-group"
	return &icon
}

func (a *ListContactsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of contacts to return (default: 10, max: 1000)").
			SetRequired(false).
			Build(),

		"offset": autoform.NewNumberField().
			SetDisplayName("Offset").
			SetDescription("Number of contacts to skip (for pagination)").
			SetRequired(false).
			Build(),

		// "filter": autoform.NewShortTextField().
		// 	SetDisplayName("Filter").
		// 	SetDescription("Filter contacts (e.g., IsExcludedFromCampaigns=true)").
		// 	SetRequired(false).Build(),
	}
}

func (a *ListContactsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(ctx.Auth.Extra["api_key"], ctx.Auth.Extra["secret_key"])
	if err != nil {
		return nil, err
	}

	limit := 10
	maxLimit := 1000
	if input.Limit > 0 {
		if input.Limit > maxLimit {
			limit = maxLimit
		} else {
			limit = input.Limit
		}
	}

	offset := 0
	if input.Offset > 0 {
		offset = input.Offset
	}

	path := "/v3/REST/contact"
	queryParams := fmt.Sprintf("?Limit=%d&Offset=%d", limit, offset)

	if input.Filter != "" {
		queryParams += "&" + input.Filter
	}

	var result map[string]interface{}
	err = client.Request(http.MethodGet, path+queryParams, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *ListContactsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListContactsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"Count": "2",
		"Data": []map[string]any{
			{
				"ContactID":               "123456",
				"Email":                   "contact1@example.com",
				"Name":                    "Contact One",
				"IsExcludedFromCampaigns": false,
				"CreatedAt":               "2023-01-01T00:00:00Z",
				"DeliveredCount":          "10",
				"IsOptInPending":          false,
				"IsSpamComplaining":       false,
				"LastActivityAt":          "2023-02-01T00:00:00Z",
			},
			{
				"ContactID":               "123457",
				"Email":                   "contact2@example.com",
				"Name":                    "Contact Two",
				"IsExcludedFromCampaigns": true,
				"CreatedAt":               "2023-01-02T00:00:00Z",
				"DeliveredCount":          "5",
				"IsOptInPending":          false,
				"IsSpamComplaining":       false,
				"LastActivityAt":          "2023-02-02T00:00:00Z",
			},
		},
		"Total": "100",
	}
}

func (a *ListContactsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (a *ListContactsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &contactListDocs,
	}
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
