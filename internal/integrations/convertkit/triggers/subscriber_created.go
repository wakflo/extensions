package triggers

import (
	"context"
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type subscriberCreatedTriggerProps struct {
	Limit int `json:"limit"`
}

type SubscriberCreatedTrigger struct{}

func (t *SubscriberCreatedTrigger) Name() string {
	return "Subscriber Created"
}

func (t *SubscriberCreatedTrigger) Description() string {
	return "Triggers a workflow when a new subscriber is added to your ConvertKit account, allowing you to automate follow-up actions or sync the data with other systems."
}

func (t *SubscriberCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *SubscriberCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &subscriberCreatedDocs,
	}
}

func (t *SubscriberCreatedTrigger) Icon() *string {
	icon := "mdi:account-alert"
	return &icon
}

func (t *SubscriberCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of subscribers to retrieve (default: 50)").
			SetDefaultValue(50).
			Build(),
	}
}

func (t *SubscriberCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *SubscriberCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *SubscriberCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[subscriberCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	limit := 50
	if input.Limit > 0 {
		limit = input.Limit
	}

	lastRunTime := ctx.Metadata().LastRun
	var fromDate string
	if lastRunTime != nil {
		fromDate = lastRunTime.Format("2006-01-02")
	} else {
		fromDate = ""
	}

	path := fmt.Sprintf("/subscribers?api_secret=%s&from=%s&page=1&limit=%d",
		ctx.Auth.Extra["api-secret"], fromDate, limit)

	response, err := shared.GetConvertKitClient(path, "GET", nil)
	if err != nil {
		return nil, errors.New("error fetching subscribers")
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, err
	}

	subscribers, ok := responseMap["subscribers"]
	if !ok {
		return nil, errors.New("failed to extract subscribers from response")
	}

	subscribersArray, ok := subscribers.([]interface{})
	if !ok {
		return nil, errors.New("invalid subscribers format in response")
	}

	if len(subscribersArray) == 0 {
		return nil, nil
	}

	return map[string]interface{}{
		"subscribers":       subscribers,
		"total_subscribers": responseMap["total_subscribers"],
	}, nil
}

func (t *SubscriberCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *SubscriberCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *SubscriberCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"subscribers": []map[string]any{
			{
				"id":            "123456",
				"first_name":    "Jane",
				"email_address": "jane@example.com",
				"state":         "active",
				"created_at":    "2023-01-15T10:30:00Z",
				"fields": map[string]string{
					"company": "Acme Inc",
				},
			},
			{
				"id":            "789012",
				"first_name":    "John",
				"email_address": "john@example.com",
				"state":         "active",
				"created_at":    "2023-01-16T14:20:00Z",
				"fields": map[string]string{
					"company": "XYZ Corp",
				},
			},
		},
		"total_subscribers": "2",
	}
}

func NewSubscriberCreatedTrigger() sdk.Trigger {
	return &SubscriberCreatedTrigger{}
}
