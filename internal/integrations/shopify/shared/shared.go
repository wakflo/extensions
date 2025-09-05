package shared

import (
	"errors"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"
)

var markdown = `
1. • Login to your Shopify account
2. • Go to Settings → Apps
3. • Click on Develop apps
4. • Create an App
5. • Fill the app name
6. • Click on Configure Admin API Scopes (Select the following scopes: 'read_orders', 'write_orders', 'write_customers', 'read_customers', 'write_products', 'read_products', 'write_draft_orders', 'read_draft_orders')
7. • Click on Install app
8. • Copy the Admin Access Token`

var (
	form = smartform.NewAuthForm("shopify-auth", "Shopify API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("domain", "Domain Name(Required)").
		Required(true).
		HelpText(`You can find your shop name in the url For example, if the URL is ` + "`" + `https://example.myshopify.com/admin` + "`" + `, then your shop name is **example**.`)

	_ = form.TextField("token", "Admin Token (Required)").
		Required(true).
		HelpText(markdown)

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
