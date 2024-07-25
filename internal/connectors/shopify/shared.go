package shopify

import (
	goshopify "github.com/bold-commerce/go-shopify/v4"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	// 	"os"
)

/*var viewStyleOptions = []*sdkcore.AutoFormSchema{
	{Const: "shopify", Title: "Shopify"},
	{Const: "custom", Title: "Shopify (Custom/Private)"},
}*/

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"domain": autoform.NewShortTextField().
			SetDisplayName("Domain Name").
			SetDescription("The domain name of the shopify app. eg. xyz.myshopify.com, type in only 'xyz'").
			SetRequired(true).
			Build(),
		"token": autoform.NewShortTextField().SetDisplayName("Authentication Token").
			SetDescription("The token used to authenticate the shopify app.").
			SetRequired(true).
			Build(),
		/*// will be enabled when dropdown cab show in dialog
		// "appMode": autoform.NewSelectField().
		//	SetDisplayName("Application Mode").
		//	SetOptions(viewStyleOptions).
		//	SetRequired(true).
		//	SetDescription("The application mode of the shopify app.").
		//	Build(),*/
	}).
	Build()

var app = goshopify.App{
	ApiKey:      "",
	ApiSecret:   "",
	RedirectUrl: "",
	Scope:       "write_orders, read_orders, write_customers, read_customers, read_products, write_products, write_draft_orders, read_draft_orders",
}

var getShopifyClient = func(shopName string, accessToken string) *goshopify.Client {
	client, err := goshopify.NewClient(app, shopName, accessToken)
	if err != nil {
		return nil
	}
	return client
}

var statusFormat = []*sdkcore.AutoFormSchema{
	{Const: "active", Title: "Active"},
	{Const: "draft", Title: "Draft"},
}

var shopifyTransactionKinds = []*sdkcore.AutoFormSchema{
	{Const: "authorization", Title: "Authorization"},
	{Const: "sale", Title: "Sale"},
	{Const: "capture", Title: "Capture"},
	{Const: "void", Title: "Void"},
	{Const: "refund", Title: "Refund"},
}
