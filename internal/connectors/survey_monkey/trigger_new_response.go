package survey_monkey

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wakflo/go-sdk/autoform"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type newResponseTriggerProps struct {
	SurveyID      string `json:"survey_id"`
	CheckInterval int    `json:"check_interval"`
}

type TriggerNewSurveyResponse struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewSurveyResponse() *TriggerNewSurveyResponse {
	return &TriggerNewSurveyResponse{
		options: &sdk.TriggerInfo{
			Name:        "New SurveyMonkey Response",
			Description: "Triggers workflow when a new response is received for the specified SurveyMonkey survey",
			RequireAuth: true,
			Auth:        sharedAuth, // Use your SurveyMonkey auth configuration
			Input: map[string]*sdkcore.AutoFormSchema{
				"survey_id": getSurveysInput("Survey ID", "Select a survey", true),
				"check_interval": autoform.NewNumberField().
					SetDisplayName("Check Interval (minutes)").
					SetDescription("How often to check for new responses (in minutes).").
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

func (t *TriggerNewSurveyResponse) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing SurveyMonkey auth token")
	}

	input := sdk.InputToType[newResponseTriggerProps](ctx)
	if input.SurveyID == "" {
		return nil, errors.New("survey_id is required")
	}
	if input.CheckInterval <= 0 {
		return nil, errors.New("check_interval must be greater than 0")
	}

	// Calculate the time range to check for new responses
	now := time.Now().UTC()
	var timeMin time.Time
	if ctx.Metadata.LastRun != nil {
		timeMin = *ctx.Metadata.LastRun
	} else {
		timeMin = now.Add(-time.Duration(input.CheckInterval) * time.Minute)
	}

	url := fmt.Sprintf("%s/surveys/%s/responses/bulk?start_modified_at=%s", baseURL,
		input.SurveyID, timeMin.Format(time.RFC3339))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+ctx.Auth.AccessToken)
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

func (t *TriggerNewSurveyResponse) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewSurveyResponse) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSurveyResponse) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSurveyResponse) GetInfo() *sdk.TriggerInfo {
	return t.options
}
