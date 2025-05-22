package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getContactActionProps struct {
	ContactID string `json:"contact_id,omitempty"`
}

type GetContactAction struct{}

func (a *GetContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_contact",
		DisplayName:   "Get Contact",
		Description:   "Retrieve detailed information about a specific contact from your Mailjet account.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getContactDocs,
		SampleOutput: map[string]any{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_contact", "Get Contact")

	shared.GetContactProp("contact_id", "Contact ID", "ID of the contact to retrieve.", false, form)

	schema := form.Build()

	return schema
}

func (a *GetContactAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(authCtx.Extra["api_key"], authCtx.Extra["secret_key"])
	if err != nil {
		return nil, err
	}

	if input.ContactID == "" {
		return nil, errors.New("name is required")
	}

	var path string

	path = "/v3/REST/contact/" + input.ContactID

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

func (a *GetContactAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetContactAction() sdk.Action {
	return &GetContactAction{}
}
