package triggers

import (
	"context"
	"errors"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type eventCreatedTriggerProps struct {
	CalendarID    string `json:"calendar_id"`
	CheckInterval int    `json:"check_interval"`
}

type EventCreatedTrigger struct{}

func (t *EventCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "event_created",
		DisplayName:   "Event Created",
		Description:   "Triggered when a new event is created in your workflow automation platform, allowing you to automate actions and workflows based on the creation of a new event.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: eventCreatedDocs,
		SampleOutput: map[string]any{
			"new_events": []map[string]any{
				{
					"id":          "abc123xyz",
					"summary":     "Team Meeting",
					"description": "Weekly team sync",
					"start": map[string]string{
						"dateTime": "2025-01-15T14:00:00Z",
					},
					"end": map[string]string{
						"dateTime": "2025-01-15T15:00:00Z",
					},
				},
			},
		},
	}
}

func (t *EventCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *EventCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *EventCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("google-calendar-event-created", "Event Created")

	form.TextField("calendar_id", "calendar_id").
		Placeholder("Calendar ID").
		HelpText("The ID of the calendar to trigger on.").
		Required(true)

	form.NumberField("check_interval", "check_interval").
		Placeholder("Check Interval (minutes)").
		HelpText("How often to check for new events (in minutes).").
		DefaultValue(5).
		Required(true)

	schema := form.Build()

	return schema
}

func (t *EventCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *EventCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *EventCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[eventCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
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

	lastRun, err := ctx.GetMetadata("lastRun")
	if err == nil && lastRun != nil {
		timeMin = *lastRun.(*time.Time)
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

func NewEventCreatedTrigger() sdk.Trigger {
	return &EventCreatedTrigger{}
}
