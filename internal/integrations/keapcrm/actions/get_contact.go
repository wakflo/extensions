package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getContactActionProps struct {
	ContactID string `json:"contact_id"`
}

type GetContactAction struct{}

func (a *GetContactAction) Name() string {
	return "Get Contact"
}

func (a *GetContactAction) Description() string {
	return "Retrieve detailed information about a specific contact in Keap"
}

func (a *GetContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"contact_id": autoform.NewShortTextField().
			SetDisplayName("Contact ID").
			SetDescription("Unique identifier of the contact to retrieve").
			SetRequired(true).Build(),
	}
}

func (a *GetContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getContactActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.ContactID == "" {
		return nil, fmt.Errorf("contact ID is required")
	}
	endpoint := "/contacts/" + input.ContactID

	contact, err := shared.MakeKeapRequest(ctx.Auth.AccessToken, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (a *GetContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getContactDocs,
	}
}

func (a *GetContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":          "12345",
		"given_name":  "John",
		"family_name": "Doe",
		"email_addresses": []map[string]string{
			{
				"email": "john.doe@example.com",
				"type":  "PRIMARY",
			},
		},
		"phone_numbers": []map[string]string{
			{
				"number": "+1-555-123-4567",
				"type":   "MOBILE",
			},
		},
		"last_updated": "2023-06-15T10:30:00Z",
	}
}

func (a *GetContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetContactAction() sdk.Action {
	return &GetContactAction{}
}

func (a *GetContactAction) Icon() *string {
	icon := "mdi:account-search"
	return &icon
}
