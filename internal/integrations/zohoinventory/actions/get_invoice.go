package actions

import (
	// "errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getInvoiceActionProps struct {
	OrganizationID string `json:"organization_id"`
	InvoiceID      string `json:"invoice_id"`
}

type GetInvoiceAction struct{}

// func (a *GetInvoiceAction) Name() string {
// 	return "Get Invoice"
// }

func (a *GetInvoiceAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_invoice",
		DisplayName:   "Get Invoice",
		Description:   "Retrieves an invoice from the accounting system, allowing you to automate tasks that require access to invoice data.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getInvoiceDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

// func (a *GetInvoiceAction) Description() string {
// 	return "Retrieves an invoice from the accounting system, allowing you to automate tasks that require access to invoice data."
// }

// func (a *GetInvoiceAction) GetType() sdkcore.ActionType {
// 	return sdkcore.ActionTypeNormal
// }

// func (a *GetInvoiceAction) Documentation() *sdk.OperationDocumentation {
// 	return &sdk.OperationDocumentation{
// 		Documentation: &getInvoiceDocs,
// 	}
// }

// func (a *GetInvoiceAction) Icon() *string {
// 	return nil
// }

// func (a *GetInvoiceAction) Properties() map[string]*sdkcore.AutoFormSchema {
// 	return map[string]*sdkcore.AutoFormSchema{
// 		"organization_id": shared.GetOrganizationsInput(),
// 		"invoice_id": autoform.NewShortTextField().
// 			SetDisplayName("Invoice ID").
// 			SetDescription("The ID of the invoice to retrieve to retrieve").
// 			SetRequired(true).
// 			Build(),
// 	}
// }

func (a *GetInvoiceAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_invoice", "Get Invoice")

	shared.GetOrganizationsProp(form)

	form.TextField("invoice_id", "Invoice ID").
		Placeholder("Enter a value for Invoice ID.").
		Required(true).
		HelpText("The ID of the invoice to retrieve to retrieve")

	schema := form.Build()

	return schema
}

func (a *GetInvoiceAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getInvoiceActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// if ctx.Auth.AccessToken == "" {
	// 	return nil, errors.New("missing Zoho auth token")
	// }

	url := fmt.Sprintf("/v1/invoices/%s?organization_id=%s",
		input.InvoiceID, input.OrganizationID)

	invoice, err := shared.GetZohoClient(authCtx.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (a *GetInvoiceAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

// func (a *GetInvoiceAction) SampleData() sdkcore.JSON {
// 	return map[string]any{
// 		"message": "Hello World!",
// 	}
// }

// func (a *GetInvoiceAction) Settings() sdkcore.ActionSettings {
// 	return sdkcore.ActionSettings{}
// }

func NewGetInvoiceAction() sdk.Action {
	return &GetInvoiceAction{}
}
