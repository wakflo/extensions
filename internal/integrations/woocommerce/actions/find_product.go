package actions

import (
	"errors"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type findProductActionProps struct {
	ProductID string `json:"productId"`
}

type FindProductAction struct{}

func (a *FindProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_product",
		DisplayName:   "Find Product",
		Description:   "Searches for a product by its ID",
		Type:          core.ActionTypeAction,
		Documentation: findProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *FindProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_product", "Find Product")

	shared.GetProductsProp("productId", "Product ID", "Enter product ID", true, form)

	schema := form.Build()

	return schema
}

func (a *FindProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[findProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	// convert ID from string to int
	productId, err := strconv.Atoi(input.ProductID)
	if err != nil {
		return nil, errors.New("Error converting to int:")
	}

	product, err := wooClient.Services.Product.One(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *FindProductAction) Auth() *core.AuthMetadata {
	return nil
}

func NewFindProductAction() sdk.Action {
	return &FindProductAction{}
}
