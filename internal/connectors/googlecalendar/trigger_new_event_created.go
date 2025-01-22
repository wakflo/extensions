package googlecalendar

import (
	"context"
	"errors"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type newEventTriggerProps struct {
	CalendarID    string `json:"calendar_id"`
	CheckInterval int    `json:"check_interval"`
}

type TriggerNewEventCreated struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewEventCreated() *TriggerNewEventCreated {
	return &TriggerNewEventCreated{
		options: &sdk.TriggerInfo{
			Name:        "New Event Created",
			Description: "Triggers workflow when a new event is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
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
			},
			Strategy: sdkcore.TriggerStrategyPolling,
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewEventCreated) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[newEventTriggerProps](ctx)
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

	if ctx.Metadata.LastRun != nil {
		timeMin = *ctx.Metadata.LastRun
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

func (t *TriggerNewEventCreated) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewEventCreated) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewEventCreated) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewEventCreated) GetInfo() *sdk.TriggerInfo {
	return t.options
}
