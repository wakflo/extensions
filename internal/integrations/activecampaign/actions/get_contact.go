package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getContactActionProps struct {
	ContactID string `json:"contact-id"`
}

type GetContactAction struct{}

func (a *GetContactAction) Name() string {
	return "Get Contact"
}

func (a *GetContactAction) Description() string {
	return "Retrieve a specific contact by ID from your ActiveCampaign account."
}

func (a *GetContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getContactDocs,
	}
}

func (a *GetContactAction) Icon() *string {
	icon := "mdi:account-search"
	return &icon
}

func (a *GetContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"contact-id": autoform.NewShortTextField().
			SetDisplayName("Contact ID").
			SetDescription("The ID of the contact you want to retrieve").
			SetRequired(true).
			Build(),
	}
}

func (a *GetContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.ContactID == "" {
		return nil, errors.New("contact ID is required")
	}

	endpoint := "contacts/" + input.ContactID

	response, err := shared.GetActiveCampaignClient(
		ctx.Auth.Extra["api_url"],
		ctx.Auth.Extra["api_key"],
		endpoint,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *GetContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":        "123",
		"email":     "sample@example.com",
		"firstName": "John",
		"lastName":  "Doe",
		"phone":     "+1234567890",
		"cdate":     "2023-01-15T15:30:00-05:00",
		"udate":     "2023-02-20T10:15:00-05:00",
		"links": map[string]string{
			"lists": "https://api.example.com/contacts/123/lists",
			"deals": "https://api.example.com/contacts/123/deals",
		},
	}
}

func (a *GetContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetContactAction() sdk.Action {
	return &GetContactAction{}
}
