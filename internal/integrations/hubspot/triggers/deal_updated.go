package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type dealUpdatedTriggerProps struct {
	Properties string `json:"properties"`
}

type DealUpdatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *DealUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "deal_updated",
		DisplayName:   "Deal Updated",
		Description:   "Trigger a workflow when deals are updated in your HubSpot CRM",
		Type:          core.TriggerTypePolling,
		Documentation: dealUpdatedDoc,
		Icon:          "mdi:cash-multiple",
		SampleOutput: map[string]any{
			"results": []map[string]any{
				{
					"id": "12345",
					"properties": map[string]any{
						"dealname":            "Enterprise Software Sale",
						"amount":              "50000",
						"dealstage":           "presentationscheduled",
						"pipeline":            "default",
						"closedate":           "2023-12-31",
						"hs_lastmodifieddate": "2023-04-15T14:30:00Z",
					},
					"createdAt": "2023-01-15T09:30:00Z",
					"updatedAt": "2023-04-15T14:30:00Z",
				},
			},
		},
	}
}

func (t *DealUpdatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

func (t *DealUpdatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

func (t *DealUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("deal_updated", "Deal Updated")

	form.TextareaField("properties", "Deal Properties").
		Required(false).
		HelpText("Comma-separated list of properties to retrieve (e.g., dealname,amount,dealstage)")

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context
func (t *DealUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations
func (t *DealUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *DealUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	props, err := sdk.InputToTypeSafely[dealUpdatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Get the last run time
	var lastRunTime *time.Time
	lr, err := ctx.GetMetadata("lastRun")
	if err == nil && lr != nil {
		lastRunTime = lr.(*time.Time)
	}

	url := "/crm/v3/objects/deals/search"
	const limit = 100

	requestBody := map[string]interface{}{
		"limit": limit,
		"sorts": []map[string]string{
			{
				"propertyName": "hs_lastmodifieddate",
				"direction":    "DESCENDING",
			},
		},
	}

	if lastRunTime != nil {
		requestBody["filterGroups"] = []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "hs_lastmodifieddate",
						"operator":     "GT",
						"value":        lastRunTime.UnixMilli(),
					},
				},
			},
		}
	}

	if props.Properties != "" {
		requestBody["properties"] = append(
			[]string{"dealname", "amount", "dealstage", "hs_lastmodifieddate"},
			props.Properties,
		)
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := shared.HubspotClient(url, authCtx.Token.AccessToken, http.MethodPost, jsonBody)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Criteria returns the criteria for triggering this trigger
func (t *DealUpdatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

// SampleData returns sample data for this trigger
func (t *DealUpdatedTrigger) SampleData() core.JSON {
	return map[string]any{
		"results": []map[string]any{
			{
				"id": "12345",
				"properties": map[string]any{
					"dealname":            "Enterprise Software Sale",
					"amount":              "50000",
					"dealstage":           "presentationscheduled",
					"pipeline":            "default",
					"closedate":           "2023-12-31",
					"hs_lastmodifieddate": "2023-04-15T14:30:00Z",
				},
				"createdAt": "2023-01-15T09:30:00Z",
				"updatedAt": "2023-04-15T14:30:00Z",
			},
		},
	}
}

func NewDealUpdatedTrigger() sdk.Trigger {
	return &DealUpdatedTrigger{}
}
