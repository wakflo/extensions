package triggers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type contactUpdatedTriggerProps struct {
	ContactsList            int   `json:"contacts_list_id,omitempty"`
	IsExcludedFromCampaigns *bool `json:"is_excluded_from_campaigns,omitempty"`
	Limit                   int   `json:"limit,omitempty"`
	LookbackHours           int   `json:"lookback_hours,omitempty"`
}

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "contact_updated",
		DisplayName:   "Contact Updated",
		Description:   "Triggers a workflow when a contact is created or updated in your Mailjet account.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: contactUpdatedDocs,
		SampleOutput: map[string]any{
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
		},
	}
}

func (t *ContactUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("contact_updated", "Contact Updated")

	form.NumberField("limit", "Limit").
		Placeholder("50").
		Required(false).
		HelpText("Maximum number of contacts to process per poll (default: 50, max: 1000)")

	schema := form.Build()

	return schema
}

func (t *ContactUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[contactUpdatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	lr, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	lastRunTime := lr.(*time.Time)

	client, err := shared.GetMailJetClient(authCtx.Extra["api_key"], authCtx.Extra["secret_key"])
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

func (t *ContactUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
