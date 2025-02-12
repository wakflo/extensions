package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/harvest/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listInvoicesActionProps struct {
	Name string `json:"name"`
}

type ListInvoicesAction struct{}

func (a *ListInvoicesAction) Name() string {
	return "List Invoices"
}

func (a *ListInvoicesAction) Description() string {
	return "Retrieve and list all invoices associated with a specific account or organization, allowing you to easily track and manage your financial transactions."
}

func (a *ListInvoicesAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListInvoicesAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listInvoicesDocs,
	}
}

func (a *ListInvoicesAction) Icon() *string {
	return nil
}

func (a *ListInvoicesAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (a *ListInvoicesAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	url := "/v2/invoices"

	invoices, err := shared.GetHarvestClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	invoiceArray, ok := invoices["invoices"].(interface{})
	if !ok {
		return nil, errors.New("failed to extract issues from response")
	}
	return invoiceArray, nil
}

func (a *ListInvoicesAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListInvoicesAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListInvoicesAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListInvoicesAction() sdk.Action {
	return &ListInvoicesAction{}
}
