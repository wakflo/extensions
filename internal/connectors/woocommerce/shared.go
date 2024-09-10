package woocommerce

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/config"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"shop-url": autoform.NewShortTextField().SetDisplayName("ShopUrl (Required)").
			SetDescription("The base URL of your app (e.g https://mystore.com) and it should start with HTTPS only").
			Build(),
		"consumer-key": autoform.NewShortTextField().SetDisplayName("Consumer Key (Required)").
			SetDescription("The consumer key generated from your app.").
			Build(),
		"consumer-secret": autoform.NewShortTextField().SetDisplayName("Consumer Secret (Required)").
			SetDescription("The consumer secret generated from your app.").
			Build(),
	}).
	Build()

func initializeWooCommerceClient(baseURL, consumerKey, consumerSecret string) *woocommerce.WooCommerce {
	const defaultTimeout = 10

	con := config.Config{
		Debug:                  true,
		URL:                    baseURL,
		Version:                "v3",
		ConsumerKey:            consumerKey,
		ConsumerSecret:         consumerSecret,
		AddAuthenticationToURL: true,
		Timeout:                defaultTimeout,
		VerifySSL:              true,
	}

	client := woocommerce.NewClient(con)
	return client
}

var productType = []*sdkcore.AutoFormSchema{
	{Const: "simple", Title: "Simple"},
	{Const: "grouped", Title: "Grouped"},
	{Const: "external", Title: "External"},
	{Const: "variable", Title: "Variable"},
}
