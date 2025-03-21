package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createContactActionProps struct {
	Email                   string `json:"email"`
	Name                    string `json:"name,omitempty"`
	IsExcludedFromCampaigns bool   `json:"is_excluded_from_campaigns,omitempty"`
}

type CreateContactAction struct{}

func (a *CreateContactAction) Name() string {
	return "Create Contact"
}

func (a *CreateContactAction) Description() string {
	return "Create a new contact in your Mailjet account."
}

func (a *CreateContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateContactAction) Icon() *string {
	icon := "mdi:account-plus"
	return &icon
}

func (a *CreateContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email address of the contact").
			SetRequired(true).Build(),

		"name": autoform.NewShortTextField().
			SetDisplayName("Name").
			SetDescription("Name of the contact").
			SetRequired(false).Build(),

		"is_excluded_from_campaigns": autoform.NewBooleanField().
			SetDisplayName("Exclude from Campaigns").
			SetDescription("Whether to exclude this contact from campaigns").
			SetRequired(false).
			SetDefaultValue(false).Build(),
	}
}

func (a *CreateContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(ctx.Auth.Extra["api_key"], ctx.Auth.Extra["secret_key"])
	if err != nil {
		return nil, err
	}

	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	payload := map[string]interface{}{
		"Email": input.Email,
	}

	if input.Name != "" {
		payload["Name"] = input.Name
	}

	payload["IsExcludedFromCampaigns"] = input.IsExcludedFromCampaigns

	var result map[string]interface{}
	err = client.Request(http.MethodPost, "/v3/REST/contact", payload, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *CreateContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"Count": 1,
		"Data": []map[string]any{
			{
				"ContactID":               "123456",
				"Email":                   "newcontact@example.com",
				"Name":                    "New Contact",
				"IsExcludedFromCampaigns": false,
				"CreatedAt":               "2023-01-01T00:00:00Z",
				"DeliveredCount":          "0",
				"IsOptInPending":          false,
				"IsSpamComplaining":       false,
				"LastActivityAt":          "2023-01-01T00:00:00Z",
			},
		},
		"Total": 1,
	}
}

func (a *CreateContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (a *CreateContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createContactDocs,
	}
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}
