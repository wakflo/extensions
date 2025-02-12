package actions

import (
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createInvoiceActionProps struct {
	Name string `json:"name"`
}

type CreateInvoiceAction struct{}

func (a *CreateInvoiceAction) Name() string {
	return "Create Invoice"
}

func (a *CreateInvoiceAction) Description() string {
	return "Create Invoice: Automatically generates and sends professional-looking invoices to customers based on predefined templates and payment terms, streamlining your accounting process and ensuring timely payments."
}

func (a *CreateInvoiceAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateInvoiceAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createInvoiceDocs,
	}
}

func (a *CreateInvoiceAction) Icon() *string {
	return nil
}

func (a *CreateInvoiceAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *CreateInvoiceAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createInvoiceActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *CreateInvoiceAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateInvoiceAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateInvoiceAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateInvoiceAction() sdk.Action {
	return &CreateInvoiceAction{}
}
