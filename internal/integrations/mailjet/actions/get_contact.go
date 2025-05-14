package actions

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getContactActionProps struct {
	ContactID int    `json:"contact_id,omitempty"`
	Email     string `json:"email,omitempty"`
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

	form.TextField("contact_id", "Contact ID").
		Placeholder("contact ID").
		Required(false).
		HelpText("ID of the contact to retrieve.")

	form.TextField("email", "Email").
		Placeholder("email").
		Required(false).
		HelpText("Email address of the contact to retrieve.")

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

func (a *GetContactAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetContactAction() sdk.Action {
	return &GetContactAction{}
}
