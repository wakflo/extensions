package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type editProductActionProps struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	BarCode     string `json:"barcode"`
	Cost        string `json:"cost"`
	ProductCode string `json:"product_code"`
	Tags        string `json:"tags"`
}

type EditProductAction struct{}

func (a *EditProductAction) Name() string {
	return "Edit Product"
}

func (a *EditProductAction) Description() string {
	return "Edit Product: Update product information by modifying existing product details, such as name, description, price, and inventory levels."
}

func (a *EditProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *EditProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &editProductDocs,
	}
}

func (a *EditProductAction) Icon() *string {
	return nil
}

func (a *EditProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Product ID").
			SetDescription("ID of product").
			SetRequired(true).Build(),
		"name": autoform.NewShortTextField().
			SetDisplayName("Product Name").
			SetDescription("name of product").
			SetRequired(false).Build(),
		"brand": autoform.NewShortTextField().
			SetDisplayName("Brand").
			SetDescription("Brand name").
			SetRequired(false).Build(),
		"category": autoform.NewShortTextField().
			SetDisplayName("Category").
			SetDescription("Category name").
			SetRequired(false).Build(),
		"product_code": autoform.NewShortTextField().
			SetDisplayName("Product Code").
			SetDescription("Product code").
			SetRequired(false).Build(),
		"barcode": autoform.NewShortTextField().
			SetDisplayName("Bar Code").
			SetDescription("Bar code").
			SetRequired(false).Build(),
		"cost": autoform.NewShortTextField().
			SetDisplayName("Product Cost").
			SetDescription("Product cost").
			SetRequired(false).Build(),
		"tags": autoform.NewLongTextField().
			SetDisplayName("Product Tags").
			SetDescription("Product tags").
			SetRequired(false).Build(),
	}
}

func (a *EditProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[editProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	formData := map[string]string{}

	if input.Name != "" {
		formData["name"] = input.Name
	}

	if input.Brand != "" {
		formData["brand"] = input.Brand
	}

	if input.Category != "" {
		formData["category"] = input.Category
	}

	if input.ProductCode != "" {
		formData["product_code"] = input.ProductCode
	}

	if input.BarCode != "" {
		formData["barcode"] = input.BarCode
	}

	if input.Cost != "" {
		formData["cost"] = input.Cost
	}

	if input.Tags != "" {
		formData["tags"] = input.Tags
	}

	endpoint := "/api/v2/edit/product/id/" + input.ID
	resp, err := shared.PrisyncRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], endpoint, http.MethodPost, formData)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *EditProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *EditProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *EditProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewEditProductAction() sdk.Action {
	return &EditProductAction{}
}
