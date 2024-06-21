package shopify

import (
	"os"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"domain": autoform.NewShortTextField().SetDisplayName("Domain Name").
			SetDescription("The domain name of the shopify app. eg. xyz.myshopify.com, type in only 'xyz' ").
			Build(),
			"token": autoform.NewShortTextField().SetDisplayName("Authentication Token").
			SetDescription("The token used to authenticate the shopify app.").
			SetRequired(true).
			Build(),
	}).
	Build();

	var app = goshopify.App{
		ApiKey:      os.Getenv("API_KEY"),
		ApiSecret:   os.Getenv("API_SECRET"),
		RedirectUrl: os.Getenv("REDIRECT_URL"),
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

