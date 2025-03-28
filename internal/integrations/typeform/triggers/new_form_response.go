package triggers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wakflo/extensions/internal/integrations/typeform/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newResponseTriggerProps struct {
	FormID        string `json:"form_id"`
	CheckInterval int    `json:"check_interval"`
}

type NewResponseTrigger struct{}

func (t *NewResponseTrigger) Name() string {
	return "New Form Response"
}

func (t *NewResponseTrigger) Description() string {
	return "Triggers workflow when a new response is received for the specified Typeform form."
}

func (t *NewResponseTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewResponseTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newFormResponseDocs,
	}
}

func (a *NewResponseTrigger) Icon() *string {
	icon := "grommet-icons:trigger"
	return &icon
}

func (t *NewResponseTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"form_id": shared.GetTypeformFormsInput("Form ID", "Select a form ID", true),
		"check_interval": autoform.NewNumberField().
			SetDisplayName("Check Interval (minutes)").
			SetDescription("How often to check for new responses (in minutes).").
			SetRequired(true).
			SetDefaultValue(5).
			Build(),
	}
}

// Start initializes the newPageCreatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewResponseTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newPageCreatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewResponseTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *NewResponseTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newResponseTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.FormID == "" {
		return nil, errors.New("form_id is required")
	}
	if input.CheckInterval <= 0 {
		return nil, errors.New("check_interval must be greater than 0")
	}

	// Calculate the time range to check for new responses
	now := time.Now().UTC()
	var timeMin time.Time
	if ctx.Metadata().LastRun != nil {
		timeMin = *ctx.Metadata().LastRun
	} else {
		timeMin = now.Add(-time.Duration(input.CheckInterval) * time.Minute)
	}

	url := fmt.Sprintf("%sforms/%s/responses?since=%s", "https://api.typeform.com/", input.FormID, timeMin.Format(time.RFC3339))

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

func (t *NewResponseTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewResponseTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewResponseTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewResponseTrigger() sdk.Trigger {
	return &NewResponseTrigger{}
}
