package actions

import (
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	sdkcore2 "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listContactsActionProps struct {
	Limit  int    `json:"limit,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type ListContactsAction struct{}

func (a *ListContactsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_contacts",
		DisplayName:   "List Contacts",
		Description:   "Retrieve a list of contacts from your Mailjet account.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: contactListDocs,
		SampleOutput: map[string]any{
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
		},
	}
}

func (a *ListContactsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_contacts", "List Contacts")

	form.NumberField("limit", "Limit").
		Placeholder("10").
		Required(false).
		HelpText("Maximum number of contacts to return (default: 10, max: 1000)")

	form.CheckboxField("filter", "Filter").
		Required(false).
		HelpText("Filter contacts (e.g., IsExcludedFromCampaigns=true)")

	schema := form.Build()

	return schema
}

func (a *ListContactsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	sdkcore2.PrettyPrint(authCtx)

	client, err := shared.GetMailJetClient(authCtx.Extra["api_key"], authCtx.Extra["secret_key"])
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

	path := "/v3/REST/contact"
	queryParams := fmt.Sprintf("?Limit=%d", limit)

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

func (a *ListContactsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
