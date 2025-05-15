package triggers

import (
	"context"
	"fmt"
	"time"

	"github.com/hiscaler/woocommerce-go"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type newProductTriggerProps struct{}

type NewProductTrigger struct{}

func (t *NewProductTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_product",
		DisplayName:   "New Product",
		Description:   "Triggered when a new product is created in your WooCommerce store.",
		Type:          core.TriggerTypePolling,
		Documentation: newProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
	}
}

func (t *NewProductTrigger) Auth() *core.AuthMetadata {
	return nil
}

func (t *NewProductTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new-product", "New Product")

	schema := form.Build()
	return schema
}

func (t *NewProductTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewProductTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewProductTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	_, err := sdk.InputToTypeSafely[newProductTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	var formattedTime string
	if lastRun != nil {
		lastRunTime, ok := lastRun.(*time.Time)
		if ok && lastRunTime != nil {
			utcTime := lastRunTime.UTC()
			formattedTime = utcTime.Format(time.RFC3339)
		}
	}

	params := woocommerce.ProductsQueryParams{
		After: formattedTime,
	}

	newProduct, total, totalPages, isLastPage, err := wooClient.Services.Product.All(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(totalPages, total, isLastPage)

	return newProduct, nil
}

func (t *NewProductTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewNewProductTrigger() sdk.Trigger {
	return &NewProductTrigger{}
}
