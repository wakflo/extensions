package shared

import (
	"errors"

	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/config"
	"github.com/wakflo/go-sdk/sdk"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var SharedAuth = autoform.NewCustomAuthField().
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

func InitializeWooCommerceClient(baseURL, consumerKey, consumerSecret string) *woocommerce.WooCommerce {
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

var ProductType = []*sdkcore.AutoFormSchema{
	{Const: "simple", Title: "Simple"},
	{Const: "grouped", Title: "Grouped"},
	{Const: "external", Title: "External"},
	{Const: "variable", Title: "Variable"},
}

func InitClient(ctx sdk.BaseContext) (*woocommerce.WooCommerce, error) {
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]
	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	wooClient := InitializeWooCommerceClient(baseURL, consumerKey, consumerSecret)
	return wooClient, nil
}
