package woocommerce

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Woocommerce",
		Description: "interacts with Woocommerce API",
		Logo:        "devicon:woocommerce",
		Version:     "0.0.1",
		Category:    sdk.Apps,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers: []sdk.ITrigger{
			NewTriggerNewProduct(),
			NewTriggerNewOrder(),
		},
		Operations: []sdk.IOperation{
			NewListProductsOperation(),
			NewFindProductsOperation(),
			NewFindCouponOperation(),
			NewUpdateProductOperation(),
			NewListOrdersOperation(),
			NewCreateProductOperation(),
			NewGetCustomerOperation(),
			NewCreateCustomerOperation(),
			NewFindCustomerOperation(),
		},
	})
}
