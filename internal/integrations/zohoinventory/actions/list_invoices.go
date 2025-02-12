package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listInvoicesActionProps struct {
	OrganizationID string `json:"organization_id"`
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
	return map[string]*sdkcore.AutoFormSchema{
		"organization_id": shared.GetOrganizationsInput(),
	}
}

func (a *ListInvoicesAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listInvoicesActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Token == nil {
		return nil, errors.New("missing Zoho auth token")
	}

	endpoint := "/v1/invoices/?organization_id=" + input.OrganizationID
	result, err := shared.GetZohoClient(ctx.Auth.Token.AccessToken, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting invoice list: %v", err)
	}

	return result, nil
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
