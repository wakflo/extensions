package zohoinventory

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getPaymentListOperationProps struct {
	OrganizationID string `json:"organization_id"`
}

type GetPaymentListOperation struct {
	options *sdk.OperationInfo
}

func NewGetPaymentListOperation() sdk.IOperation {
	return &GetPaymentListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Customer Payment List",
			Description: "Retrieve a list of Payments",
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

func (c *GetPaymentListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getPaymentListOperationProps](ctx)

	url := "https://www.zohoapis.com/inventory/v1/customerpayments/?organization_id=" + input.OrganizationID

	paymentList, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return paymentList, nil
}

func (c *GetPaymentListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetPaymentListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
