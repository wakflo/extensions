package actions

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type emailInvoiceActionProps struct {
	Name string `json:"name"`
}

type EmailInvoiceAction struct{}

func (a *EmailInvoiceAction) Name() string {
	return "Email Invoice"
}

func (a *EmailInvoiceAction) Description() string {
	return "Sends an email to the customer with a detailed invoice summary, including payment instructions and any relevant notes or attachments."
}

func (a *EmailInvoiceAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *EmailInvoiceAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &emailInvoiceDocs,
	}
}

func (a *EmailInvoiceAction) Icon() *string {
	return nil
}

func (a *EmailInvoiceAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *EmailInvoiceAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[emailInvoiceActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *EmailInvoiceAction) Auth() *sdk.Auth {
	return nil
}

func (a *EmailInvoiceAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *EmailInvoiceAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewEmailInvoiceAction() sdk.Action {
	return &EmailInvoiceAction{}
}
