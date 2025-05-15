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

type newOrderTriggerProps struct {
	Tag string `json:"tag"`
}

type NewOrderTrigger struct{}

func (t *NewOrderTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_order",
		DisplayName:   "New Order",
		Description:   "Triggered when a new order is created in your WooCommerce store.",
		Type:          core.TriggerTypePolling,
		Documentation: newOrderDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
	}
}

func (t *NewOrderTrigger) Auth() *core.AuthMetadata {
	return nil
}

func (t *NewOrderTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new-order", "New Order")

	schema := form.Build()
	return schema
}

func (t *NewOrderTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewOrderTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewOrderTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	_, err := sdk.InputToTypeSafely[newOrderTriggerProps](ctx)
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

	params := woocommerce.OrdersQueryParams{
		After: formattedTime,
	}

	newOrder, total, totalPages, isLastPage, err := wooClient.Services.Order.All(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(totalPages, total, isLastPage)

	return newOrder, nil
}

func (t *NewOrderTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewNewOrderTrigger() sdk.Trigger {
	return &NewOrderTrigger{}
}
