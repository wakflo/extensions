package triggers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type contactUpdatedTriggerProps struct {
	ContactsList            int   `json:"contacts_list_id,omitempty"`
	IsExcludedFromCampaigns *bool `json:"is_excluded_from_campaigns,omitempty"`
	Limit                   int   `json:"limit,omitempty"`
	LookbackHours           int   `json:"lookback_hours,omitempty"`
}

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Name() string {
	return "Contact Updated"
}

func (t *ContactUpdatedTrigger) Description() string {
	return "Triggers a workflow when a contact is created or updated in your Mailjet account."
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
	icon := "mdi:account-edit"
	return &icon
}

func (t *ContactUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of contacts to process per poll (default: 50, max: 1000)").
			SetRequired(false).
			Build(),
	}
}

func (t *ContactUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	input, err := sdk.InputToTypeSafely[contactUpdatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(ctx.Auth.Extra["api_key"], ctx.Auth.Extra["secret_key"])
	if err != nil {
		return nil, err
	}

	limit := 50
	if input.Limit > 0 && input.Limit <= 1000 {
		limit = input.Limit
	}

	url := "/v3/REST/contactdata"
	queryParams := fmt.Sprintf("?Limit=%d", limit)

	var result map[string]interface{}
	err = client.Request(http.MethodGet, url+queryParams, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated contacts: %v", err)
	}

	data, ok := result["Data"].([]interface{})
	if !ok || len(data) == 0 {
		return []interface{}{}, nil
	}

	if lastRunTime == nil {
		return result, nil
	}

	filteredData := make([]interface{}, 0)
	for _, item := range data {
		contact, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		if lastActivity, ok := contact["LastActivityAt"].(string); ok {
			lastActivityTime, err := time.Parse(time.RFC3339, lastActivity)
			if err == nil && lastActivityTime.After(*lastRunTime) {
				filteredData = append(filteredData, contact)
			}
		}
	}

	if len(filteredData) == 0 {
		return []interface{}{}, nil
	}

	filteredResult := map[string]interface{}{
		"Count": len(filteredData),
		"Data":  filteredData,
		"Total": len(filteredData),
	}

	return filteredResult, nil
}

func (t *ContactUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"Count": "2",
		"Data": []map[string]any{
			{
				"ContactID":               "123456",
				"Email":                   "contact1@example.com",
				"Name":                    "Contact One Updated",
				"IsExcludedFromCampaigns": false,
				"CreatedAt":               "2023-01-01T00:00:00Z",
				"LastActivityAt":          "2023-04-01T14:25:36Z",
				"DeliveredCount":          "15",
				"IsOptInPending":          false,
				"IsSpamComplaining":       false,
			},
			{
				"ContactID":               "123457",
				"Email":                   "newcontact@example.com",
				"Name":                    "New Contact",
				"IsExcludedFromCampaigns": false,
				"CreatedAt":               "2023-04-01T12:30:00Z",
				"LastActivityAt":          "2023-04-01T12:30:00Z",
				"DeliveredCount":          "0",
				"IsOptInPending":          true,
				"IsSpamComplaining":       false,
			},
		},
		"Total": "2",
	}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
