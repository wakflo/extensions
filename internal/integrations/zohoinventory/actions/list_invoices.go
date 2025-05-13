package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listInvoicesActionProps struct {
	OrganizationID string `json:"organization_id"`
}

type ListInvoicesAction struct{}

// func (a *ListInvoicesAction) Name() string {
// 	return "List Invoices"
// }

func (a *ListInvoicesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_invoices",
		DisplayName:   "List Invoices",
		Description:   "Retrieve and list all invoices associated with a specific account or organization, allowing you to easily track and manage your financial transactions.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listInvoicesDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

// func (a *ListInvoicesAction) Description() string {
// 	return "Retrieve and list all invoices associated with a specific account or organization, allowing you to easily track and manage your financial transactions."
// }

// func (a *ListInvoicesAction) GetType() sdkcore.ActionType {
// 	return sdkcore.ActionTypeNormal
// }

// func (a *ListInvoicesAction) Documentation() *sdk.OperationDocumentation {
// 	return &sdk.OperationDocumentation{
// 		Documentation: &listInvoicesDocs,
// 	}
// }

// func (a *ListInvoicesAction) Icon() *string {
// 	return nil
// }

// func (a *ListInvoicesAction) Properties() map[string]*sdkcore.AutoFormSchema {
// 	return map[string]*sdkcore.AutoFormSchema{
// 		"organization_id": shared.GetOrganizationsInput(),
// 	}
// }

func (a *ListInvoicesAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("list_invoices", "Lisrt Invoices")

	shared.GetOrganizationsProp(form)

	schema := form.Build()

	return schema
}

func (a *ListInvoicesAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listInvoicesActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// if ctx.Auth.Token == nil {
	// 	return nil, errors.New("missing Zoho auth token")
	// }

	endpoint := "/v1/invoices/?organization_id=" + input.OrganizationID
	result, err := shared.GetZohoClient(authCtx.AccessToken, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting invoice list: %v", err)
	}

	return result, nil
}

func (a *ListInvoicesAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

// func (a *ListInvoicesAction) SampleData() sdkcore.JSON {
// 	return map[string]any{
// 		"message": "Hello World!",
// 	}
// }

// func (a *ListInvoicesAction) Settings() sdkcore.ActionSettings {
// 	return sdkcore.ActionSettings{}
// }

func NewListInvoicesAction() sdk.Action {
	return &ListInvoicesAction{}
}
