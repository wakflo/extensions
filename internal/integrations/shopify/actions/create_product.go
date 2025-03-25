package actions

import (
	"context"
	"errors"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/shopspring/decimal"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *CreateProductAction) Name() string {
	return "Create Product"
}

func (a *CreateProductAction) Description() string {
	return "Create Product: Automatically generates and creates new products in your system, including product details such as name, description, price, and inventory levels."
}

func (a *CreateProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createProductDocs,
	}
}

func (a *CreateProductAction) Icon() *string {
	return nil
}

func (a *CreateProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"title": autoform.NewShortTextField().
			SetDisplayName("Product title").
			SetDescription("The title of the product.").
			SetRequired(true).
			Build(),
		"bodyHTML": autoform.NewLongTextField().
			SetDisplayName("Product description").
			SetDescription("The description of the product.").
			SetRequired(false).
			Build(),
		"vendor": autoform.NewShortTextField().
			SetDisplayName("Vendor").
			SetDescription("Vendor.").
			SetRequired(false).
			Build(),
		"productType": autoform.NewShortTextField().
			SetDisplayName("Product type").
			SetDescription("A categorization for the product used for filtering and searching products.").
			SetRequired(false).
			Build(),
		"tags": autoform.NewLongTextField().
			SetDisplayName("Tags").
			SetDescription("A string of comma-separated tags for filtering and search.").
			SetRequired(false).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Product status").
			SetDescription("The status of the product: active or draft").
			SetOptions(shared.StatusFormat).
			SetRequired(false).
			Build(),
		"imageURL": autoform.NewShortTextField().
			SetDisplayName("Image URL").
			SetDescription("URL for the product image.").
			SetRequired(false).
			Build(),
		"price": autoform.NewShortTextField().
			SetDisplayName("Price").
			SetDescription("The price of the product.").
			SetRequired(false).
			Build(),
	}
}

func (a *CreateProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
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

func (a *CreateProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateProductAction() sdk.Action {
	return &CreateProductAction{}
}
