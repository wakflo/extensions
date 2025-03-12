package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createNewContactActionProps struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	MobileNumber *string `json:"mobile_number"`
	Email        *string `json:"email"`
	JobTitle     *string `json:"job_title,omitempty"`
	Company      *string `json:"company,omitempty"`
	Address      *string `json:"address,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	ZipCode      *string `json:"zip_code,omitempty"`
	Country      *string `json:"country,omitempty"`
}

type CreateNewContactAction struct{}

func (a *CreateNewContactAction) Name() string {
	return "Create a Contact"
}

func (a *CreateNewContactAction) Description() string {
	return "Create a new contact in Freshworks CRM with specified information."
}

func (a *CreateNewContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateNewContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createContactDocs,
	}
}

func (a *CreateNewContactAction) Icon() *string {
	return nil
}

func (a *CreateNewContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"first_name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("Contact's first name").
			SetRequired(true).
			Build(),
		"last_name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Contact's last name").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Contact's email").
			SetRequired(true).
			Build(),
		"mobile_number": autoform.NewShortTextField().
			SetDisplayName("Mobile Number").
			SetDescription("Contact's mobile number").
			SetRequired(false).
			Build(),
		"job_title": autoform.NewShortTextField().
			SetDisplayName("Job Title").
			SetDescription("Job title of the contact").
			SetRequired(false).Build(),
		"company": autoform.NewShortTextField().
			SetDisplayName("Company").
			SetDescription("Company name of the contact").
			SetRequired(false).Build(),
		"address": autoform.NewShortTextField().
			SetDisplayName("Address").
			SetDescription("Street address of the contact").
			SetRequired(false).Build(),
		"city": autoform.NewShortTextField().
			SetDisplayName("City").
			SetDescription("City of the contact").
			SetRequired(false).Build(),
		"state": autoform.NewShortTextField().
			SetDisplayName("State").
			SetDescription("State/Province of the contact").
			SetRequired(false).Build(),
		"zip_code": autoform.NewShortTextField().
			SetDisplayName("Zip Code").
			SetDescription("Postal/Zip code of the contact").
			SetRequired(false).Build(),
		"country": autoform.NewShortTextField().
			SetDisplayName("Country").
			SetDescription("Country of the contact").
			SetRequired(false).Build(),
	}
}

func (a *CreateNewContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	input := sdk.InputToType[createNewContactActionProps](ctx.BaseContext)

	domain := ctx.Auth.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	contactData := map[string]interface{}{
		"contact": map[string]interface{}{
			"first_name":    input.FirstName,
			"last_name":     input.LastName,
			"email":         input.Email,
			"mobile_number": input.MobileNumber,
		},
	}

	contactMap := contactData["contact"].(map[string]interface{})
	if input.MobileNumber != nil {
		contactMap["phone"] = *input.MobileNumber
	}
	if input.JobTitle != nil {
		contactMap["job_title"] = *input.JobTitle
	}
	if input.Company != nil {
		contactMap["company"] = *input.Company
	}
	if input.Address != nil {
		contactMap["address"] = *input.Address
	}
	if input.City != nil {
		contactMap["city"] = *input.City
	}
	if input.State != nil {
		contactMap["state"] = *input.State
	}
	if input.ZipCode != nil {
		contactMap["zipcode"] = *input.ZipCode
	}
	if input.Country != nil {
		contactMap["country"] = *input.Country
	}

	response, err := shared.CreateContact(freshworksDomain, ctx.Auth.Extra["api-key"], contactData)
	if err != nil {
		return nil, fmt.Errorf("error creating contact:  %v", err)
	}

	return response, nil
}

func (a *CreateNewContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateNewContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"contact": map[string]any{
			"id":         "12345",
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john.doe@example.com",
			"created_at": "2023-01-01T12:00:00Z",
			"updated_at": "2023-01-01T12:00:00Z",
		},
	}
}

func (a *CreateNewContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateNewContactAction() sdk.Action {
	return &CreateNewContactAction{}
}
