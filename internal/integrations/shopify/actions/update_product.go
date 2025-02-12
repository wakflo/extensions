package actions

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateProductActionProps struct {
	ProductID   uint64 `json:"productID"`
	Title       string `json:"title"`
	BodyHTML    string `json:"bodyHTML"`
	Vendor      string `json:"vendor"`
	ProductType string `json:"productType"`
	Status      string `json:"status"`
	Tags        string `json:"tags"`
}

type UpdateProductAction struct{}

func (a *UpdateProductAction) Name() string {
	return "Update Product"
}

func (a *UpdateProductAction) Description() string {
	return "Updates product information in your e-commerce platform or CRM system by mapping to specific fields such as product name, description, price, and inventory levels."
}

func (a *UpdateProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateProductDocs,
	}
}

func (a *UpdateProductAction) Icon() *string {
	return nil
}

func (a *UpdateProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"productID": autoform.NewNumberField().
			SetDisplayName("Product ID").
			SetDescription("The id of the product.").
			SetRequired(true).
			Build(),
		"title": autoform.NewShortTextField().
			SetDisplayName("Product title").
			SetDescription("The title of the product.").
			SetRequired(false).
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
			SetDisplayName("Status").
			SetDescription("The status of the product: active or draft").
			SetOptions(shared.StatusFormat).
			SetRequired(true).
			Build(),
	}
}

func (a *UpdateProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	existingProduct, err := client.Product.Get(context.Background(), input.ProductID, nil)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		existingProduct.Title = input.Title
	}
	if input.BodyHTML != "" {
		existingProduct.BodyHTML = input.BodyHTML
	}
	if input.Vendor != "" {
		existingProduct.Vendor = input.Vendor
	}
	if input.ProductType != "" {
		existingProduct.ProductType = input.ProductType
	}
	if input.Tags != "" {
		existingProduct.Tags = input.Tags
	}
	if input.Status != "" {
		existingProduct.Status = goshopify.ProductStatus(input.Status)
	}

	updatedProduct, err := client.Product.Update(context.Background(), *existingProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return map[string]interface{}{
		"updated_product": updatedProduct,
	}, nil
}

func (a *UpdateProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateProductAction() sdk.Action {
	return &UpdateProductAction{}
}
