package zohoinventory

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getInvoiceOperationProps struct {
	OrganizationID string `json:"organization_id"`
	InvoiceID      string `json:"invoice_id"`
}

type GetInvoiceOperation struct {
	options *sdk.OperationInfo
}

func NewGetInvoiceOperation() sdk.IOperation {
	return &GetInvoiceOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Invoice",
			Description: "Retrieve a specific invoice.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": autoform.NewShortTextField().
					SetDisplayName("Organization ID").
					SetDescription("The Zoho Inventory organization ID").
					SetRequired(true).
					Build(),
				"invoice_id": autoform.NewShortTextField().
					SetDisplayName("Payment ID").
					SetDescription("The ID of the customer payment to retrieve").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetInvoiceOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getInvoiceOperationProps](ctx)

	url := fmt.Sprintf("https://www.zohoapis.com/inventory/v1/invoices/%s?organization_id=%s",
		input.InvoiceID, input.OrganizationID)
	invoice, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (c *GetInvoiceOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetInvoiceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
