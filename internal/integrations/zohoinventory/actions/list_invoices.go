package actions

import (
	"errors"
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
	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	endpoint := "/v1/invoices/?organization_id=" + input.OrganizationID
	result, err := shared.GetZohoClient(token, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting invoice list: %v", err)
	}

	return result, nil
}

func (a *ListInvoicesAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListInvoicesAction() sdk.Action {
	return &ListInvoicesAction{}
}
