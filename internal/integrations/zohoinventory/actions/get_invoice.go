package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getInvoiceActionProps struct {
	OrganizationID string `json:"organization_id"`
	InvoiceID      string `json:"invoice_id"`
}

type GetInvoiceAction struct{}

func (a *GetInvoiceAction) Name() string {
	return "Get Invoice"
}

func (a *GetInvoiceAction) Description() string {
	return "Retrieves an invoice from the accounting system, allowing you to automate tasks that require access to invoice data."
}

func (a *GetInvoiceAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetInvoiceAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getInvoiceDocs,
	}
}

func (a *GetInvoiceAction) Icon() *string {
	return nil
}

func (a *GetInvoiceAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"organization_id": shared.GetOrganizationsInput(),
		"invoice_id": autoform.NewShortTextField().
			SetDisplayName("Invoice ID").
			SetDescription("The ID of the invoice to retrieve to retrieve").
			SetRequired(true).
			Build(),
	}
}

func (a *GetInvoiceAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getInvoiceActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Zoho auth token")
	}

	url := fmt.Sprintf("/v1/invoices/%s?organization_id=%s",
		input.InvoiceID, input.OrganizationID)

	invoice, err := shared.GetZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (a *GetInvoiceAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetInvoiceAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetInvoiceAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetInvoiceAction() sdk.Action {
	return &GetInvoiceAction{}
}
