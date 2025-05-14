package actions

import (
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type createContactActionProps struct {
	Email                   string `json:"email"`
	Name                    string `json:"name,omitempty"`
	IsExcludedFromCampaigns bool   `json:"is_excluded_from_campaigns,omitempty"`
}

type CreateContactAction struct{}

func (a *CreateContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_contact",
		DisplayName:   "Create Contact",
		Description:   "Create a new contact in your Mailjet account.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createContactDocs,
		SampleOutput: map[string]any{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateContactAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("create_contact", "Create Contact")

	form.TextField("email", "Email").
		Placeholder("example@example.com").
		Required(true).
		HelpText("Email address of the contact")

	form.TextField("name", "Name").
		Placeholder("John Doe").
		Required(false).
		HelpText("Name of the contact")

	form.CheckboxField("is_excluded_from_campaigns", "Exclude from Campaigns").
		Placeholder("Enter a value for Exclude from Campaigns.").
		Required(false).
		DefaultValue(false).
		HelpText("Whether to exclude this contact from campaigns")

	schema := form.Build()

	return schema
}

func (a *CreateContactAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createContactActionProps](ctx)
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

func (a *CreateContactAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateContactAction() sdk.Action {
	return &CreateContactAction{}
}
