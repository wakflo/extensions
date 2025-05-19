package actions

import (
	// "errors"
	"errors"
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
	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	url := fmt.Sprintf("/v1/invoices/%s?organization_id=%s",
		input.InvoiceID, input.OrganizationID)

	invoice, err := shared.GetZohoClient(token, url)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (a *GetInvoiceAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetInvoiceAction() sdk.Action {
	return &GetInvoiceAction{}
}
