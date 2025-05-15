package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listOrdersActionProps struct {
	Limit int `json:"limit"`
}

type ListOrdersAction struct{}

func (a *ListOrdersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_orders",
		DisplayName:   "List Orders",
		Description:   "Retrieve a list of orders from yoour WooCommerce store.",
		Type:          core.ActionTypeAction,
		Documentation: listOrdersDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListOrdersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_orders", "List Orders")

	form.NumberField("limit", "Result Limit").
		Placeholder("Maximum number of orders to return").
		HelpText("Maximum number of orders to return")

	schema := form.Build()

	return schema
}

func (a *ListOrdersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listOrdersActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	params := woocommerce.OrdersQueryParams{}

	if input.Limit > 0 {
		params.PerPage = input.Limit
	}

	orders, _, _, _, err := wooClient.Services.Order.All(params)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (a *ListOrdersAction) Auth() *core.AuthMetadata {
	return nil
}

func NewListOrdersAction() sdk.Action {
	return &ListOrdersAction{}
}
