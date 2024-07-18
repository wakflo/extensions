package woocommerce

import (
	"errors"
	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/config"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type ListProductsOperation struct {
	options *sdk.OperationInfo
}

func NewListProductsOperation() *ListProductsOperation {
	return &ListProductsOperation{
		options: &sdk.OperationInfo{
			Name:        "List products",
			Description: "List products in store",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId": autoform.NewLongTextField().
					SetDisplayName("").
					SetDescription("").
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListProductsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	//URL := "https://shut-present.localsite.io/wp-json/wc/v3"

	con := config.Config{
		Debug:                  true,
		URL:                    "https://shut-present.localsite.io/",
		Version:                "v3",
		ConsumerKey:            consumerKey,
		ConsumerSecret:         consumerSecret,
		AddAuthenticationToURL: false,
		Timeout:                10,
		VerifySSL:              true,
	}

	wooClient := woocommerce.NewClient(con)

	// Create a query parameters struct
	//params := woocommerce.ProductsQueryParams{}

	// Get all products
	product, _ := wooClient.Services.Product.One(1)

	// Process the products
	//for _, product := range products {
	//	fmt.Printf("Product ID: %d, Name: %s, Price: %s\n", product.ID, product.Name, product.Price)
	//}

	return product, nil
}

func (c *ListProductsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListProductsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
