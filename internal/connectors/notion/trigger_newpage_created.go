package notion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type newPageTriggerProps struct {
	DatabaseID    string `json:"database"`
	CheckInterval int    `json:"check_interval"`
}

type TriggerNewPageCreated struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewPageCreated() *TriggerNewPageCreated {
	return &TriggerNewPageCreated{
		options: &sdk.TriggerInfo{
			Name:        "New Notion Page Created",
			Description: "Triggers workflow when a new page is created in the specified Notion database",
			RequireAuth: true,
			Auth:        sharedAuth, // Define sharedAuth based on how your application handles Notion auth tokens
			Input: map[string]*sdkcore.AutoFormSchema{
				"database": getNotionDatabasesInput("Database", "Select a database", true),
				"check_interval": autoform.NewNumberField().
					SetDisplayName("Check Interval (minutes)").
					SetDescription("How often to check for new pages (in minutes).").
					SetRequired(true).
					SetDefaultValue(5).
					Build(),
			},
			Type:     sdkcore.TriggerTypeCron,
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewPageCreated) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing Notion auth token")
	}

	input := sdk.InputToType[newPageTriggerProps](ctx)

	if input.DatabaseID == "" {
		return nil, errors.New("database_id is required")
	}

	if input.CheckInterval <= 0 {
		return nil, errors.New("check_interval must be greater than 0")
	}

	now := time.Now().UTC()
	var timeMin time.Time

	if ctx.Metadata.LastRun != nil {
		timeMin = *ctx.Metadata.LastRun
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

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", input.DatabaseID), bytes.NewBuffer(filterJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+ctx.Auth.Token.AccessToken)
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
				"title":      getPageTitle(page.Properties),
				"properties": page.Properties,
			})
		}
		return map[string]interface{}{"new_pages": newPages}, nil
	}

	return map[string]interface{}{"message": "No new pages found"}, nil
}

func (t *TriggerNewPageCreated) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewPageCreated) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewPageCreated) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewPageCreated) GetInfo() *sdk.TriggerInfo {
	return t.options
}

// / Helper function to get the title from page properties
func getPageTitle(properties map[string]interface{}) string {
	if titleProp, ok := properties["Name"].(map[string]interface{}); ok {
		if titleArray, ok := titleProp["title"].([]interface{}); ok && len(titleArray) > 0 {
			if titleText, ok := titleArray[0].(map[string]interface{}); ok {
				if plainText, ok := titleText["plain_text"].(string); ok {
					return plainText
				}
			}
		}
	}
	return "Untitled"
}
