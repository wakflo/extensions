package triggers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/notion/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type newPageCreatedTriggerProps struct {
	DatabaseID    string `json:"database"`
	CheckInterval int    `json:"check_interval"`
}

type NewPageCreatedTrigger struct{}

func (t *NewPageCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_page_created",
		DisplayName:   "New Page Created",
		Description:   "Triggers when a new page is created in your website or application, allowing you to automate tasks and workflows immediately after page creation.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newPageCreatedDocs,
		SampleOutput: map[string]any{
			"page": map[string]any{
				"id": "1234567890",
			},
		},
	}
}

func (t *NewPageCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewPageCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new_page_created", "New Page Created")

	shared.GetNotionDatabasesProp(form)
	form.NumberField("check_interval", "Check Interval (minutes)").
		Required(true).
		HelpText("How often to check for new pages (in minutes).").
		DefaultValue(5)

	schema := form.Build()

	return schema
}

// Start initializes the newPageCreatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewPageCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newPageCreatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewPageCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newPageCreatedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewPageCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newPageCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	if input.DatabaseID == "" {
		return nil, errors.New("database_id is required")
	}

	if input.CheckInterval <= 0 {
		return nil, errors.New("check_interval must be greater than 0")
	}

	now := time.Now().UTC()
	var timeMin time.Time

	lr, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	if lr != nil {
		timeMin = lr.(*time.Time).UTC()
	} else {
		timeMin = now.Add(-time.Duration(input.CheckInterval) * time.Minute)
	}

	filter := map[string]interface{}{
		"filter": map[string]interface{}{
			"timestamp": "created_time",
			"created_time": map[string]interface{}{
				"on_or_after": timeMin.Format(time.RFC3339),
			},
		},
	}

	filterJSON, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode filter JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(shared.BaseURL+"/databases/%s/query", input.DatabaseID), bytes.NewBuffer(filterJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notion API returned non-OK status: %s", resp.Status)
	}

	var queryResult struct {
		Results []struct {
			ID          string                 `json:"id"`
			CreatedTime time.Time              `json:"created_time"`
			Properties  map[string]interface{} `json:"properties"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&queryResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(queryResult.Results) > 0 {
		newPages := make([]map[string]interface{}, 0, len(queryResult.Results))
		for _, page := range queryResult.Results {
			newPages = append(newPages, map[string]interface{}{
				"id":         page.ID,
				"created_at": page.CreatedTime,
				"title":      shared.GetPageTitle(page.Properties),
				"properties": page.Properties,
			})
		}
		return map[string]interface{}{"new_pages": newPages}, nil
	}

	return map[string]interface{}{"message": "No new pages found"}, nil
}

func (t *NewPageCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewPageCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewPageCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewPageCreatedTrigger() sdk.Trigger {
	return &NewPageCreatedTrigger{}
}
