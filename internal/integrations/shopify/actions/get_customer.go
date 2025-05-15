package actions

import (
	"context"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getCustomerActionProps struct {
	CustomerID uint64 `json:"customerId"`
}

type GetCustomerAction struct{}

func (a *GetCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_customer",
		DisplayName:   "Get Customer",
		Description:   "Retrieves customer information from a specified data source or system, allowing you to access and utilize existing customer data within your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: getCustomerDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_customer", "Get Customer")

	form.NumberField("customerId", "Customer ID").
		Placeholder("The ID of the customer.").
		Required(true).
		HelpText("The ID of the customer.")

	schema := form.Build()
	return schema
}

func (a *GetCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCustomerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	customer, err := client.Customer.Get(context.Background(), input.CustomerID, nil)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, fmt.Errorf("no customer found with ID '%d'", input.CustomerID)
	}

	return core.JSON(map[string]interface{}{
		"raw_customer": customer,
	}), nil
}

func (a *GetCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetCustomerAction() sdk.Action {
	return &GetCustomerAction{}
}
