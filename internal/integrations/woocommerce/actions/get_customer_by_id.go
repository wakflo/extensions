package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getCustomerByIDActionProps struct {
	CustomerID int `json:"customer-id"`
}

type GetCustomerByIDAction struct{}

func (a *GetCustomerByIDAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_customer_by_id",
		DisplayName:   "Get Customer By ID",
		Description:   "Retrieves a customer record by their unique identifier (ID)",
		Type:          core.ActionTypeAction,
		Documentation: getCustomerByIDDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetCustomerByIDAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_customer_by_id", "Get Customer By ID")

	form.NumberField("customer-id", "Customer ID").
		Placeholder("Enter customer ID").
		Required(true).
		HelpText("Enter customer ID")

	schema := form.Build()

	return schema
}

func (a *GetCustomerByIDAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCustomerByIDActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	customer, err := wooClient.Services.Customer.One(input.CustomerID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (a *GetCustomerByIDAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetCustomerByIDAction() sdk.Action {
	return &GetCustomerByIDAction{}
}
