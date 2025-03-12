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

type contactUpdatedTriggerProps struct {
	Properties string `json:"properties"`
}

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Name() string {
	return "Contact Updated"
}

func (t *ContactUpdatedTrigger) Description() string {
	return "Trigger a workflow when a contact is created or updated in your HubSpot CRM"
}

func (t *ContactUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &contactUpdatedDoc,
	}
}

func (t *ContactUpdatedTrigger) Icon() *string {
	icon := "mdi:account-check"
	return &icon
}

func (t *ContactUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"properties": autoform.NewLongTextField().
			SetDisplayName("Contact Properties").
			SetDescription("Comma-separated list of contact properties to include in the response (e.g., firstname,lastname,email)").
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
	input, err := sdk.InputToTypeSafely[contactUpdatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun

	url := "/crm/v3/objects/contacts/search"
	const limit = 100

	requestBody := map[string]interface{}{
		"limit": limit,
		"sorts": []map[string]string{
			{
				"propertyName": "lastmodifieddate",
				"direction":    "DESCENDING",
			},
		},
	}

	if lastRunTime != nil {
		requestBody["filterGroups"] = []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "lastmodifieddate",
						"operator":     "GT",
						"value":        lastRunTime.UnixMilli(),
					},
				},
			},
		}
	}

	if input.Properties != "" {
		requestBody["properties"] = append(
			[]string{"firstname", "lastname", "email", "lastmodifieddate"},
			input.Properties,
		)
	}

	if input.Properties != "" {
		requestBody["properties"] = append(
			[]string{"firstname", "lastname", "email", "lastmodifieddate"},
			input.Properties,
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

func (t *ContactUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"results": []map[string]any{
			{
				"id": "51",
				"properties": map[string]any{
					"firstname":        "John",
					"lastname":         "Doe",
					"email":            "john.doe@example.com",
					"phone":            "+1234567890",
					"createdate":       "2023-03-15T09:31:40.678Z",
					"lastmodifieddate": "2023-03-15T10:45:12.412Z",
				},
				"createdAt": "2023-03-15T09:31:40.678Z",
				"updatedAt": "2023-03-15T10:45:12.412Z",
			},
		},
	}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
