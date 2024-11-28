package typeform

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
	FormID        string `json:"form_id"`
	CheckInterval int    `json:"check_interval"`
}

type TriggerNewTypeformResponse struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewFormResponse() *TriggerNewTypeformResponse {
	return &TriggerNewTypeformResponse{
		options: &sdk.TriggerInfo{
			Name:        "New Typeform Response",
			Description: "Triggers workflow when a new response is received for the specified Typeform form",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"form_id": getTypeformFormsInput("Form ID", "Select a form ID", true),
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

func (t *TriggerNewTypeformResponse) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing Typeform auth token")
	}

	input := sdk.InputToType[newResponseTriggerProps](ctx)
	if input.FormID == "" {
		return nil, errors.New("form_id is required")
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

	url := fmt.Sprintf("%sforms/%s/responses?since=%s", baseURL, input.FormID, timeMin.Format(time.RFC3339))

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
		return nil, fmt.Errorf("typeform API returned non-OK status: %s", resp.Status)
	}

	// Parse Typeform response
	var triggerResponseResults struct {
		Items []struct {
			Token       string                   `json:"token"`
			SubmittedAt time.Time                `json:"submitted_at"`
			Answers     []map[string]interface{} `json:"answers"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&triggerResponseResults); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(triggerResponseResults.Items) > 0 {
		newResponses := make([]map[string]interface{}, 0, len(triggerResponseResults.Items))
		for _, response := range triggerResponseResults.Items {
			newResponses = append(newResponses, map[string]interface{}{
				"token":        response.Token,
				"submitted_at": response.SubmittedAt,
				"answers":      response.Answers,
			})
		}
		return map[string]interface{}{"new_responses": newResponses}, nil
	}

	return map[string]interface{}{"message": "No new responses found"}, nil
}

func (t *TriggerNewTypeformResponse) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewTypeformResponse) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewTypeformResponse) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewTypeformResponse) GetInfo() *sdk.TriggerInfo {
	return t.options
}
