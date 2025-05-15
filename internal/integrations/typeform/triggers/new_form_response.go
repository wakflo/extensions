package triggers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/typeform/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type newResponseTriggerProps struct {
	FormID        string `json:"form_id"`
	CheckInterval int    `json:"check_interval"`
}

type NewResponseTrigger struct{}

func (t *NewResponseTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_form_response",
		DisplayName:   "New Form Response",
		Description:   "Triggers workflow when a new response is received for the specified Typeform form.",
		Type:          core.TriggerTypePolling,
		Documentation: newFormResponseDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Icon: "grommet-icons:trigger",
	}
}

func (t *NewResponseTrigger) Auth() *core.AuthMetadata {
	return nil
}

func (t *NewResponseTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new-form-response", "New Form Response")

	shared.RegisterTypeformFormsProps(form, "Form ID", "Select a form ID", true)

	form.NumberField("check_interval", "Check Interval (minutes)").
		Placeholder("How often to check for new responses (in minutes).").
		Required(true).
		DefaultValue(5).
		HelpText("How often to check for new responses (in minutes).")

	schema := form.Build()
	return schema
}

// Start initializes the newResponseTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewResponseTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newResponseTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewResponseTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *NewResponseTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[newResponseTriggerProps](ctx)
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

	// Get the last run time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	if lastRun != nil {
		lastRunTime, ok := lastRun.(*time.Time)
		if ok && lastRunTime != nil {
			timeMin = *lastRunTime
		} else {
			timeMin = now.Add(-time.Duration(input.CheckInterval) * time.Minute)
		}
	} else {
		timeMin = now.Add(-time.Duration(input.CheckInterval) * time.Minute)
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%sforms/%s/responses?since=%s", "https://api.typeform.com/", input.FormID, timeMin.Format(time.RFC3339))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
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

func (t *NewResponseTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewNewResponseTrigger() sdk.Trigger {
	return &NewResponseTrigger{}
}
