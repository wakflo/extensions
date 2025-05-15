package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type subscriberCreatedTriggerProps struct {
	Limit int `json:"limit"`
}

type SubscriberCreatedTrigger struct{}

func (t *SubscriberCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "subscriber_created",
		DisplayName:   "Subscriber Created",
		Description:   "Triggers a workflow when a new subscriber is added to your ConvertKit account, allowing you to automate follow-up actions or sync the data with other systems.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: subscriberCreatedDocs,
		Icon:          "mdi:account-alert",
		SampleOutput: map[string]any{
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
		},
	}
}

func (t *SubscriberCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *SubscriberCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *SubscriberCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("convertkit-subscriber-created", "Subscriber Created")

	form.NumberField("limit", "limit").
		Placeholder("Limit").
		HelpText("Maximum number of subscribers to retrieve (default: 50)").
		DefaultValue(50).
		Required(false)

	schema := form.Build()

	return schema
}

func (t *SubscriberCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *SubscriberCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *SubscriberCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[subscriberCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	limit := 50
	if input.Limit > 0 {
		limit = input.Limit
	}

	var lastRunTime *time.Time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err == nil && lastRun != nil {
		lastRunTime = lastRun.(*time.Time)
	}

	var fromDate string
	if lastRunTime != nil {
		fromDate = lastRunTime.Format("2006-01-02")
	} else {
		fromDate = ""
	}

	path := fmt.Sprintf("/subscribers?api_secret=%s&from=%s&page=1&limit=%d",
		authCtx.Extra["api-secret"], fromDate, limit)

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

func NewSubscriberCreatedTrigger() sdk.Trigger {
	return &SubscriberCreatedTrigger{}
}
