package actions

import (
	"errors"
	"net/http"

	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getContactActionProps struct {
	ContactID string `json:"contact_id"`
}

type GetContactAction struct{}

func (a *GetContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Get Contact",
		Description:   "Retrieve detailed information about a specific contact in Keap",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getContactDocs,
		SampleOutput: map[string]any{
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
		},
	}
}

func (a *GetContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_contact", "Get Contact")

	form.TextField("contact_id", "Contact ID").
		Placeholder("contact ID").
		Required(true).
		HelpText("Unique identifier of the contact to retrieve")

	schema := form.Build()

	return schema
}

func (a *GetContactAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	if input.ContactID == "" {
		return nil, errors.New("contact ID is required")
	}
	endpoint := "/contacts/" + input.ContactID

	contact, err := shared.MakeKeapRequest(token, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (a *GetContactAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetContactAction() sdk.Action {
	return &GetContactAction{}
}
