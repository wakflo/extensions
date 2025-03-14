package triggers

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type ContactCreatedTrigger struct{}

func (t *ContactCreatedTrigger) Name() string {
	return "Contact Created"
}

func (t *ContactCreatedTrigger) Description() string {
	return "Trigger a workflow when a new contact is created in Keap"
}

func (t *ContactCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &contactCreatedDocs,
	}
}

func (t *ContactCreatedTrigger) Icon() *string {
	icon := "mdi:account-plus-outline"
	return &icon
}

func (t *ContactCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (t *ContactCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun
	var createdSince string
	if lastRunTime != nil {
		createdSince = lastRunTime.UTC().Format(time.RFC3339Nano)
	} else {
		createdSince = ""
	}

	queryParams := url.Values{}
	queryParams.Add("limit", "100") // Maximum allowed limit
	queryParams.Add("order", "date_created")
	queryParams.Add("order_direction", "ascending")
	queryParams.Add("since", createdSince)

	queryParams.Add("optional_properties", "email")
	queryParams.Add("optional_properties", "last_updated")

	endpoint := fmt.Sprintf("/contacts?%s", queryParams.Encode())

	response, err := shared.MakeKeapRequest(ctx.Auth.AccessToken, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (t *ContactCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"contacts": []map[string]any{
			{
				"id":          "123",
				"given_name":  "John",
				"family_name": "Doe",
				"email_addresses": []map[string]string{
					{
						"email": "john.doe@example.com",
						"type":  "PRIMARY",
					},
				},
				"date_created": "2023-06-15T14:30:25Z",
				"last_updated": "2023-06-15T14:30:25Z",
			},
		},
	}
}

func NewContactCreatedTrigger() sdk.Trigger {
	return &ContactCreatedTrigger{}
}
