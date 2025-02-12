package actions

import (
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getInvoiceActionProps struct {
	TenantID  string `json:"tenant_id"`
	InvoiceID string `json:"invoice_id"`
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
		"tenant_id":  shared.GetTenantInput("Organization", "select organization", true),
		"invoice_id": shared.GetInvoiceInput("invoices", "invoices", false),
	}
}

func (a *GetInvoiceAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getInvoiceActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	endpoint := "/Invoices/" + input.InvoiceID

	invoice, err := shared.GetXeroNewClient(ctx.Auth.AccessToken, endpoint, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoice: %v", err)
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
