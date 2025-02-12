package actions

import (
	"context"
	"errors"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createProductActionProps struct {
	Title       string `json:"title"`
	BodyHTML    string `json:"bodyHTML"`
	Vendor      string `json:"vendor"`
	ProductType string `json:"productType"`
	Status      string `json:"status"`
	Tags        string `json:"tags"`
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
	product, err := client.Product.Create(context.Background(), newProduct)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not created ")
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
