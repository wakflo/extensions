package triggers

import (
	"context"
	"errors"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type eventCreatedTriggerProps struct {
	CalendarID    string `json:"calendar_id"`
	CheckInterval int    `json:"check_interval"`
}

type EventCreatedTrigger struct{}

func (t *EventCreatedTrigger) Name() string {
	return "Event Created"
}

func (t *EventCreatedTrigger) Description() string {
	return "Triggered when a new event is created in your workflow automation platform, allowing you to automate actions and workflows based on the creation of a new event."
}

func (t *EventCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *EventCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &eventCreatedDocs,
	}
}

func (t *EventCreatedTrigger) Icon() *string {
	return nil
}

func (t *EventCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"calendar_id": autoform.NewShortTextField().
			SetDisplayName("Calendar ID").
			SetDescription("The ID of the calendar to trigger on.").
			SetRequired(true).
			Build(),
		"check_interval": autoform.NewNumberField().
			SetDisplayName("Check Interval (minutes)").
			SetDescription("How often to check for new events (in minutes).").
			SetRequired(true).
			SetDefaultValue(5).
			Build(),
	}
}

// Start initializes the eventCreatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *EventCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the eventCreatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *EventCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of eventCreatedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *EventCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[eventCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.CalendarID == "" {
		return nil, errors.New("calendar_id is required")
	}

	if input.CheckInterval <= 0 {
		return nil, errors.New("check_interval must be greater than 0")
	}

	now := time.Now().UTC()
	var timeMin time.Time

	if ctx.Metadata().LastRun != nil {
		timeMin = *ctx.Metadata().LastRun
	} else {
		timeMin = now.Add(-time.Duration(input.CheckInterval) * time.Minute)
	}

	events, err := srv.Events.List(input.CalendarID).
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(now.Format(time.RFC3339)).
		OrderBy("startTime").
		SingleEvents(true).
		Do()
	if err != nil {
		return nil, err
	}

	if len(events.Items) > 0 {
		newEvents := make([]map[string]interface{}, 0, len(events.Items))
		for _, event := range events.Items {
			newEvents = append(newEvents, map[string]interface{}{
				"id":          event.Id,
				"summary":     event.Summary,
				"description": event.Description,
				"start":       event.Start,
				"end":         event.End,
			})
		}
		return map[string]interface{}{"new_events": newEvents}, nil
	}

	return map[string]interface{}{"message": "No new events found"}, nil
}

func (t *EventCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *EventCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *EventCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewEventCreatedTrigger() sdk.Trigger {
	return &EventCreatedTrigger{}
}
