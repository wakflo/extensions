package shared

import (
	"errors"
	"fmt"

	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/config"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("woocommerce-auth", "WooCommerce API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("shop-url", "Shop Url (Required)").
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

func GetProductsProp(id string, title, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getProducts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		wooClient, err := InitClient(ctx)
		if err != nil {
			return nil, err
		}

		params := woocommerce.ProductsQueryParams{}

		products, _, _, _, err := wooClient.Services.Product.All(params)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch WooCommerce products: %v", err)
		}

		var options []map[string]interface{}

		for _, product := range products {
			options = append(options, map[string]interface{}{
				"id":   fmt.Sprintf("%d", product.ID),
				"name": product.Name,
			})
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select product").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getProducts)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

func GetCustomersProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getCustomers := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Initialize WooCommerce client
		wooClient, err := InitClient(ctx)
		if err != nil {
			return nil, err
		}

		params := woocommerce.CustomersQueryParams{}

		// Make request to WooCommerce API using the proper service method
		customers, _, _, _, err := wooClient.Services.Customer.All(params)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch WooCommerce customers: %v", err)
		}

		// Process the response
		var options []map[string]interface{}

		// Extract customer information from response
		for _, customer := range customers {
			options = append(options, map[string]interface{}{
				"id":   fmt.Sprintf("%d", customer.ID),
				"name": fmt.Sprintf("%s %s", customer.FirstName, customer.LastName),
			})
		}

		// Return the formatted options
		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select customer").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getCustomers)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

func GetCouponsProp(id string, title, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getCoupons := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Initialize WooCommerce client
		wooClient, err := InitClient(ctx)
		if err != nil {
			return nil, err
		}

		params := woocommerce.CouponsQueryParams{}

		// Make request to WooCommerce API using the proper service method
		coupons, _, _, _, err := wooClient.Services.Coupon.All(params)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch WooCommerce coupons: %v", err)
		}

		// Process the response
		var options []map[string]interface{}

		// Extract coupon information from response
		for _, coupon := range coupons {
			options = append(options, map[string]interface{}{
				"id":   fmt.Sprintf("%d", coupon.ID),
				"name": coupon.Code,
			})
		}

		// Return the formatted options
		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select coupon").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getCoupons)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}
