package actions

import (
	"context"
	"errors"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/shopspring/decimal"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createProductActionProps struct {
	Title       string           `json:"title"`
	BodyHTML    string           `json:"bodyHTML"`
	Vendor      string           `json:"vendor"`
	ProductType string           `json:"productType"`
	Status      string           `json:"status"`
	Tags        string           `json:"tags"`
	ImageURL    string           `json:"imageURL"`
	Price       *decimal.Decimal `json:"price"`
}

type CreateProductAction struct{}

func (a *CreateProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_product",
		DisplayName:   "Create Product",
		Description:   "Create Product: Automatically generates and creates new products in your system, including product details such as name, description, price, and inventory levels.",
		Type:          core.ActionTypeAction,
		Documentation: createProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_product", "Create Product")

	form.TextField("title", "Product title").
		Placeholder("The title of the product.").
		Required(true).
		HelpText("The title of the product.")

	form.TextareaField("bodyHTML", "Product description").
		Placeholder("The description of the product.").
		HelpText("The description of the product.")

	form.TextField("vendor", "Vendor").
		Placeholder("Vendor.").
		HelpText("Vendor.")

	form.TextField("productType", "Product type").
		Placeholder("A categorization for the product used for filtering and searching products.").
		HelpText("A categorization for the product used for filtering and searching products.")

	form.TextareaField("tags", "Tags").
		Placeholder("A string of comma-separated tags for filtering and search.").
		HelpText("A string of comma-separated tags for filtering and search.")

	form.SelectField("status", "Product status").
		AddOption("active", "Active").
		AddOption("draft", "draft").
		Placeholder("The status of the product: active or draft").
		HelpText("The status of the product: active or draft")

	form.TextField("imageURL", "Image URL").
		Placeholder("URL for the product image.").
		HelpText("URL for the product image.")

	form.TextField("price", "Price").
		Placeholder("The price of the product.").
		HelpText("The price of the product.")

	schema := form.Build()
	return schema
}

func (a *CreateProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	newProduct := goshopify.Product{
		Title:       input.Title,
		BodyHTML:    input.BodyHTML,
		Vendor:      input.Vendor,
		ProductType: input.ProductType,
		Status:      goshopify.ProductStatus(input.Status),
		Tags:        input.Tags,
	}

	// Add price via variant if provided
	if input.Price != nil && input.Price.Sign() > 0 {
		// Create a default variant with the price
		variant := goshopify.Variant{
			Option1: "Default",   // Default option title
			Price:   input.Price, // Pass the decimal.Decimal pointer directly
		}
		newProduct.Variants = []goshopify.Variant{variant}
	}

	// Create the product with all details including variant/price
	product, err := client.Product.Create(context.Background(), newProduct)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not created")
	}

	if input.ImageURL != "" {
		image := goshopify.Image{
			ProductId: product.Id,
			Src:       input.ImageURL,
		}

		_, err = client.Image.Create(context.Background(), product.Id, image)
		if err != nil {
			return map[string]interface{}{
				"new product": product,
				"warning":     "Product created but failed to add image: " + err.Error(),
			}, nil
		}
	}

	return map[string]interface{}{
		"new product": product,
	}, nil
}

func (a *CreateProductAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateProductAction() sdk.Action {
	return &CreateProductAction{}
}
