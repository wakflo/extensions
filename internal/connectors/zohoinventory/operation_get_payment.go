package zohoinventory

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getPaymentOperationProps struct {
	OrganizationID string `json:"organization_id"`
	PaymentID      string `json:"payment_id"`
}

type GetPaymentOperation struct {
	options *sdk.OperationInfo
}

func NewGetPaymentOperation() sdk.IOperation {
	return &GetPaymentOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Customer Payment",
			Description: "Retrieve a specific payment.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": autoform.NewShortTextField().
					SetDisplayName("Organization ID").
					SetDescription("The Zoho Inventory organization ID").
					SetRequired(true).
					Build(),
				"payment_id": autoform.NewShortTextField().
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

func (c *GetPaymentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getPaymentOperationProps](ctx)

	url := fmt.Sprintf("https://www.zohoapis.com/inventory/v1/customerpayments/%s?organization_id=%s",
		input.PaymentID, input.OrganizationID)

	payments, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (c *GetPaymentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetPaymentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
