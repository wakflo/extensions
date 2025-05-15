package shared

import (
	"errors"

	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/config"
	"github.com/juicycleff/smartform/v1"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"
)

var (
	form = smartform.NewAuthForm("woocommerce-auth", "WooCommerce API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("shop-url", "ShopUrl (Required)").
		Required(true).
		HelpText("The base URL of your app (e.g https://mystore.com) and it should start with HTTPS only")

	_ = form.TextField("consumer-key", "Consumer Key (Required)").
		Required(true).
		HelpText("The consumer key generated from your app.")

	_ = form.TextField("consumer-secret", "Consumer Secret (Required)").
		Required(true).
		HelpText("The consumer secret generated from your app.")

	WooSharedAuth = form.Build()
)

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

func InitClient(ctx sdkcontext.BaseContext) (*woocommerce.WooCommerce, error) {
	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra == nil {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	baseURL := authCtx.Extra["shop-url"]
	consumerKey := authCtx.Extra["consumer-key"]
	consumerSecret := authCtx.Extra["consumer-secret"]
	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	wooClient := InitializeWooCommerceClient(baseURL, consumerKey, consumerSecret)
	return wooClient, nil
}
