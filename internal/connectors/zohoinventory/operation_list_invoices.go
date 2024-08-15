package zohoinventory

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getInvoiceListOperationProps struct {
	OrganizationID string `json:"organization_id"`
}

type GetInvoiceListOperation struct {
	options *sdk.OperationInfo
}

func NewGetInvoiceListOperation() sdk.IOperation {
	return &GetInvoiceListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Invoice List",
			Description: "Retrieve a list of Invoices",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": autoform.NewShortTextField().
					SetDisplayName("Organization ID").
					SetDescription("The Zoho Inventory organization ID").
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

func (c *GetInvoiceListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getInvoiceListOperationProps](ctx)

	url := "https://www.zohoapis.com/inventory/v1/invoices/?organization_id=" + input.OrganizationID

	invoices, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(invoices), nil
}

func (c *GetInvoiceListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetInvoiceListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
