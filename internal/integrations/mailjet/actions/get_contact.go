package actions

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getContactActionProps struct {
	ContactID int    `json:"contact_id,omitempty"`
	Email     string `json:"email,omitempty"`
}

type GetContactAction struct{}

func (a *GetContactAction) Name() string {
	return "Get Contact"
}

func (a *GetContactAction) Description() string {
	return "Retrieve detailed information about a specific contact from your Mailjet account."
}

func (a *GetContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetContactAction) Icon() *string {
	icon := "mdi:account-details"
	return &icon
}

func (a *GetContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"contact_id": autoform.NewNumberField().
			SetDisplayName("Contact ID").
			SetDescription("ID of the contact to retrieve").
			SetRequired(false).Build(),

		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email address of the contact to retrieve").
			SetRequired(false).Build(),
	}
}

func (a *GetContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(ctx.Auth.Extra["api_key"], ctx.Auth.Extra["secret_key"])
	if err != nil {
		return nil, err
	}

	if input.ContactID <= 0 && input.Email == "" {
		return nil, fmt.Errorf("either Contact ID or Email must be provided")
	}

	var path string
	if input.ContactID > 0 {
		path = fmt.Sprintf("/v3/REST/contact/%d", input.ContactID)
	} else {
		path = "/v3/REST/contact"
		path += "?Email=" + url.QueryEscape(input.Email)
	}

	var result map[string]interface{}
	err = client.Request(http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, err
	}

	if count, ok := result["Count"].(float64); ok && count == 0 {
		return nil, fmt.Errorf("contact not found")
	}

	return result, nil
}

func (a *GetContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"Count": 1,
		"Data": []map[string]any{
			{
				"ContactID":               "123456",
				"Email":                   "contact@example.com",
				"Name":                    "Example Contact",
				"IsExcludedFromCampaigns": false,
				"CreatedAt":               "2023-01-01T00:00:00Z",
				"DeliveredCount":          "10",
				"IsOptInPending":          false,
				"IsSpamComplaining":       false,
				"LastActivityAt":          "2023-02-01T00:00:00Z",
			},
		},
		"Total": 1,
	}
}

func (a *GetContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (a *GetContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getContactDocs,
	}
}

func NewGetContactAction() sdk.Action {
	return &GetContactAction{}
}
