package triggers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/surveymonkey/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type newResponseTriggerProps struct {
	SurveyID      string `json:"survey_id"`
	CheckInterval int    `json:"check_interval"`
}

type NewResponseTrigger struct{}

func (t *NewResponseTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_response",
		DisplayName:   "New Response",
		Description:   "Triggers a workflow when a new response is submitted to a specific survey in your SurveyMonkey account.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newResponseDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
	}
}

func (t *NewResponseTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewResponseTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new_response", "New Response")

	shared.GetSurveysProp("survey_id", "Survey ID", "The ID of the survey to monitor for new responses.", true, form)

	schema := form.Build()

	return schema
}

// Start initializes the newPageCreatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewResponseTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newPageCreatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewResponseTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewResponseTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newResponseTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	if input.SurveyID == "" {
		return nil, errors.New("form_id is required")
	}

	// Calculate the time range to check for new responses
	lr, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	lastRunTime := lr.(*time.Time)

	now := time.Now().UTC()
	var timeMin time.Time

	if lastRunTime != nil {
		timeMin = *lastRunTime
	} else {
		timeMin = now
	}

	url := fmt.Sprintf("%s/surveys/%s/responses/bulk?start_modified_at=%s", "https://api.surveymonkey.com/v3",
		input.SurveyID, timeMin.Format(time.RFC3339))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SurveyMonkey API returned non-OK status: %s", resp.Status)
	}

	// Parse SurveyMonkey response
	var responseResults struct {
		Data []struct {
			ID          string    `json:"id"`
			DateCreated time.Time `json:"date_created"`
			Answers     []struct {
				QuestionID string      `json:"question_id"`
				Answer     interface{} `json:"answer"`
			} `json:"answers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseResults); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(responseResults.Data) > 0 {
		newResponses := make([]map[string]interface{}, 0, len(responseResults.Data))
		for _, response := range responseResults.Data {
			newResponses = append(newResponses, map[string]interface{}{
				"id":           response.ID,
				"date_created": response.DateCreated,
				"answers":      response.Answers,
			})
		}
		return map[string]interface{}{"new_responses": newResponses}, nil
	}

	return map[string]interface{}{"message": "No new responses found"}, nil
}

func (t *NewResponseTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewResponseTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewNewResponseTrigger() sdk.Trigger {
	return &NewResponseTrigger{}
}
