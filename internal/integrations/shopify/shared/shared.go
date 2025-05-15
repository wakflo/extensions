package shared

import (
	"errors"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"
)

var (
	form = smartform.NewAuthForm("shopify-auth", "Shopify API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("domain", "Domain Name(Required)").
		Required(true).
		HelpText("The domain name of the shopify app. eg. xyz.myshopify.com, type in only 'xyz'")

	_ = form.TextField("token", "Consumer Key (Required)").
		Required(true).
		HelpText("The consumer key generated from your app.")

	ShopifySharedAuth = form.Build()
)

var app = goshopify.App{
	ApiKey:      "",
	ApiSecret:   "",
	RedirectUrl: "",
	Scope:       "write_orders, read_orders, write_customers, read_customers, read_products, write_products, write_draft_orders, read_draft_orders",
}

var GetShopifyClient = func(shopName string, accessToken string) *goshopify.Client {
	client, err := goshopify.NewClient(app, shopName, accessToken)
	if err != nil {
		return nil
	}
	return client
}

func CreateClient(ctx sdkcontext.BaseContext) (*goshopify.Client, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	domain := authCtx.Extra["domain"]
	shopName := domain + ".myshopify.com"

	return GetShopifyClient(shopName, authCtx.Extra["token"]), nil
}
