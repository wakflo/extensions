package triggers

import (
	"context"
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type subscriberAddedTriggerProps struct {
	ListID string `json:"listId"`
}

type SubscriberAddedTrigger struct{}

func (t *SubscriberAddedTrigger) Name() string {
	return "Subscriber Added"
}

func (t *SubscriberAddedTrigger) Description() string {
	return "Trigger a workflow when a new subscriber is added to a list."
}

func (t *SubscriberAddedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *SubscriberAddedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &subscriberAddedDocs,
	}
}

func (t *SubscriberAddedTrigger) Icon() *string {
	icon := "mdi:account-plus-outline"
	return &icon
}

func (t *SubscriberAddedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"listId": shared.GetCreateSendSubscriberListsInput(),
	}
}

func (t *SubscriberAddedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *SubscriberAddedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *SubscriberAddedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[subscriberAddedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.ListID == "" {
		return nil, errors.New("list ID is required")
	}

	lastRunTime := ctx.Metadata().LastRun

	endpoint := fmt.Sprintf("lists/%s/active.json", input.ListID)

	if lastRunTime != nil {
		dateParam := lastRunTime.Format("2006-01-02")
		endpoint += "?date=" + dateParam
	}

	response, err := shared.GetCampaignMonitorClient(
		ctx.Auth.Extra["api-key"],
		ctx.Auth.Extra["client-id"],
		endpoint,
		"GET",
		nil)
	if err != nil {
		return nil, err
	}

	results, ok := response.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format")
	}

	subscribers, ok := results["Results"]
	if !ok {
		return nil, errors.New("invalid response format: Results field not found")
	}

	var sinceDateStr string
	if lastRunTime != nil {
		sinceDateStr = lastRunTime.Format("2006-01-02")
	} else {
		sinceDateStr = "N/A"
	}

	return map[string]interface{}{
		"newSubscribers": subscribers,
		"listId":         input.ListID,
		"sinceDate":      sinceDateStr,
	}, nil
}

func (t *SubscriberAddedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *SubscriberAddedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *SubscriberAddedTrigger) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"newSubscribers": []interface{}{
			map[string]interface{}{
				"EmailAddress": "subscriber1@example.com",
				"Name":         "John Doe",
				"Date":         "2023-07-10T14:30:00",
				"State":        "Active",
				"CustomFields": []interface{}{
					map[string]interface{}{
						"Key":   "City",
						"Value": "New York",
					},
					map[string]interface{}{
						"Key":   "Age",
						"Value": "34",
					},
				},
			},
			map[string]interface{}{
				"EmailAddress": "subscriber2@example.com",
				"Name":         "Jane Smith",
				"Date":         "2023-06-20T09:15:00",
				"State":        "Active",
				"CustomFields": []interface{}{
					map[string]interface{}{
						"Key":   "City",
						"Value": "Los Angeles",
					},
					map[string]interface{}{
						"Key":   "Age",
						"Value": "29",
					},
				},
			},
		},
		"listId":    "a1b2c3d4e5f6g7h8i9j0",
		"sinceDate": "2023-07-09",
	}
}

func NewSubscriberAddedTrigger() sdk.Trigger {
	return &SubscriberAddedTrigger{}
}
