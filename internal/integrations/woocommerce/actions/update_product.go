package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateProductActionProps struct {
	ProductID        int     `json:"productId"`
	Name             string  `json:"name"`
	RegularPrice     float64 `json:"regular_price"`
	SalePrice        float64 `json:"sale_price"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Weight           string  `json:"weight"`
	Length           float64 `json:"length"`
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	ImageURL         string  `json:"image_url"`
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
		"productId": autoform.NewNumberField().
			SetDisplayName("Product ID").
			SetDescription("Enter the product ID").
			SetRequired(true).
			Build(),
		"name": autoform.NewShortTextField().
			SetDisplayName("Product Name").
			SetDescription("Enter product Name").
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName(" Description").
			SetDescription("Enter product description").
			Build(),
		"short_description": autoform.NewLongTextField().
			SetDisplayName("Short Description").
			SetDescription("Enter the short description").
			Build(),
		"length": autoform.NewNumberField().
			SetDisplayName("Length").
			SetDescription("Enter Product Length").
			Build(),
		"regular_price": autoform.NewNumberField().
			SetDisplayName("Regular Price").
			SetDescription("Enter Regular Price").
			SetRequired(true).
			Build(),
		"sale_price": autoform.NewNumberField().
			SetDisplayName("Discounted Price").
			SetDescription("Enter Discounted Price").
			SetRequired(true).
			Build(),
		"height": autoform.NewNumberField().
			SetDisplayName("Height").
			SetDescription("Enter Product Height").
			Build(),
		"width": autoform.NewNumberField().
			SetDisplayName("Width").
			SetDescription("Enter Product Width").
			Build(),
		"weight": autoform.NewShortTextField().
			SetDisplayName("Weight").
			SetDescription("weight").
			Build(),
		"image_url": autoform.NewShortTextField().
			SetDisplayName("Image URL").
			SetDescription("Enter image URL. Must end with a valid image extension (.jpg, .jpeg, .png, .gif)").
			Build(),
	}
}

func (a *UpdateProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	params := woocommerce.UpdateProductRequest{}
	if input.Name != "" {
		params.Name = input.Name
	}

	if input.Description != "" {
		params.Description = input.Description
	}

	if input.RegularPrice != 0 {
		params.RegularPrice = input.RegularPrice
	}

	if input.SalePrice != 0 {
		params.SalePrice = input.SalePrice
	}

	if input.ShortDescription != "" {
		params.ShortDescription = input.ShortDescription
	}

	if input.Weight != "" {
		params.Weight = input.Weight
	}

	if input.Length != 0 {
		params.Dimensions.Length = input.Length
	}

	if input.Width != 0 {
		params.Dimensions.Width = input.Width
	}

	if input.Height != 0 {
		params.Dimensions.Length = input.Height
	}

	if input.ImageURL != "" {
		params.Images = []entity.ProductImage{
			{
				Src: input.ImageURL,
			},
		}
	}

	updatedProduct, err := wooClient.Services.Product.Update(input.ProductID, params)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
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
