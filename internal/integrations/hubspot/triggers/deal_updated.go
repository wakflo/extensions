package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type dealUpdatedTriggerProps struct {
	Properties string `json:"properties"`
}

type DealUpdatedTrigger struct{}

func (t *DealUpdatedTrigger) Name() string {
	return "Deal Updated"
}

func (t *DealUpdatedTrigger) Description() string {
	return "Trigger a workflow when deals are updated in your HubSpot CRM"
}

func (t *DealUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *DealUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &dealUpdatedDoc,
	}
}

func (t *DealUpdatedTrigger) Icon() *string {
	icon := "mdi:cash-multiple"
	return &icon
}

func (t *DealUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"properties": autoform.NewShortTextField().
			SetDisplayName("Deal Properties").
			SetDescription("Comma-separated list of properties to retrieve (e.g., dealname,amount,dealstage)").
			SetRequired(false).
			Build(),
	}
}

func (t *DealUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *DealUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *DealUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	props, err := sdk.InputToTypeSafely[dealUpdatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun
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

	resp, err := shared.HubspotClient(url, ctx.Auth.AccessToken, http.MethodPost, jsonBody)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *DealUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *DealUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *DealUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"results": []map[string]any{},
	}
}

func NewDealUpdatedTrigger() sdk.Trigger {
	return &DealUpdatedTrigger{}
}
