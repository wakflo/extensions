package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type emailInvoiceActionProps struct {
	Name string `json:"name"`
}

type EmailInvoiceAction struct{}

func (a *EmailInvoiceAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "email_invoice",
		DisplayName:   "Email Invoice",
		Description:   "Sends an email to the customer with a detailed invoice summary, including payment instructions and any relevant notes or attachments.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: emailInvoiceDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *EmailInvoiceAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("email_invoice", "Email Invoice")

	form.TextField("name", "Name").
		Placeholder("Name").
		Required(true).
		HelpText("Your name")

	schema := form.Build()

	return schema
}

func (a *EmailInvoiceAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[emailInvoiceActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *EmailInvoiceAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewEmailInvoiceAction() sdk.Action {
	return &EmailInvoiceAction{}
}
