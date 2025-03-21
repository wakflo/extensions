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

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Name() string {
	return "Contact Updated"
}

func (t *ContactUpdatedTrigger) Description() string {
	return "Trigger a workflow when a contact is updated in Keap"
}

func (t *ContactUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &contactUpdatedDocs,
	}
}

func (t *ContactUpdatedTrigger) Icon() *string {
	icon := "mdi:account-edit-outline"
	return &icon
}

func (t *ContactUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (t *ContactUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	// Determine the time range to check for updated contacts
	lastRunTime := ctx.Metadata().LastRun
	var updatedSince string
	if lastRunTime != nil {
		updatedSince = lastRunTime.UTC().Format(time.RFC3339Nano)
	} else {
		updatedSince = ""
	}

	queryParams := url.Values{}
	queryParams.Add("limit", "100")
	queryParams.Add("order", "last_updated")
	queryParams.Add("order_direction", "ascending")
	queryParams.Add("since", updatedSince)

	endpoint := "/contacts?" + queryParams.Encode()

	response, err := shared.MakeKeapRequest(ctx.Auth.AccessToken, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (t *ContactUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactUpdatedTrigger) SampleData() sdkcore.JSON {
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
				"date_created": "2023-06-15T14:30:25.039Z",
				"last_updated": "2023-06-16T14:30:25.039Z",
			},
		},
	}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
