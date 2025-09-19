package actions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlecalendar/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type createEventActionProps struct {
	CalendarID    string `json:"calendar_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Location      string `json:"location"`
	StartDateTime string `json:"start_datetime"`
	EndDateTime   string `json:"end_datetime"`
}

type CreateEventAction struct{}

func (a *CreateEventAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_event",
		DisplayName:   "Create Event",
		Description:   "Create Event: Triggers the creation of a new event in your chosen calendar or scheduling system, allowing you to automate the process of setting up meetings, appointments, and other events from within your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: createEventDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"id":      "event_123",
			"summary": "New Meeting",
			"start": map[string]any{
				"dateTime": "2025-09-23T10:00:00Z",
			},
			"end": map[string]any{
				"dateTime": "2025-09-23T11:00:00Z",
			},
			"htmlLink": "https://calendar.google.com/event?eid=...",
			"created":  "2025-09-23T09:00:00Z",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateEventAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_event", "Create Event")

	shared.RegisterCalendarProps(form)

	form.TextField("title", "Title").
		Placeholder("Event Title").
		HelpText("The title of the event.").
		Required(true)

	form.TextareaField("description", "Description").
		Placeholder("Event Description").
		HelpText("The description of the event").
		Required(true)

	form.TextField("location", "Location").
		Placeholder("Event Location").
		HelpText("The location of the event").
		Required(false)

	form.DateField("start_datetime", "Start Date & Time").
		Placeholder("Select start date and time").
		HelpText("The start date and time of the event").
		Required(true)

	form.DateField("end_datetime", "End Date & Time").
		Placeholder("Select end date and time").
		HelpText("The end date and time of the event").
		Required(true)

	schema := form.Build()

	return schema
}

func (a *CreateEventAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *CreateEventAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createEventActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	eventService, err := calendar.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.CalendarID == "" {
		return nil, errors.New("calendar id is required")
	}

	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.Description == "" {
		return nil, errors.New("description is required")
	}

	if input.StartDateTime == "" {
		return nil, errors.New("start date and time is required")
	}

	if input.EndDateTime == "" {
		return nil, errors.New("end date and time is required")
	}

	startDateTime, err := shared.ParseDateTime(input.StartDateTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start datetime: %v", err)
	}

	endDateTime, err := shared.ParseDateTime(input.EndDateTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end datetime: %v", err)
	}

	if endDateTime.Before(startDateTime) {
		return nil, errors.New("end datetime must be after start datetime")
	}

	duration := endDateTime.Sub(startDateTime)
	if duration > 24*7*time.Hour {
		return nil, errors.New("event duration cannot exceed 7 days")
	}

	startDateTimeStr := startDateTime.Format(time.RFC3339)
	endDateTimeStr := endDateTime.Format(time.RFC3339)

	event, err := eventService.Events.Insert(input.CalendarID, &calendar.Event{
		Summary:     input.Title,
		Description: input.Description,
		Location:    input.Location,
		Start: &calendar.EventDateTime{
			DateTime: startDateTimeStr,
		},
		End: &calendar.EventDateTime{
			DateTime: endDateTimeStr,
		},
	}).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %v", err)
	}

	return event, nil
}

func NewCreateEventAction() sdk.Action {
	return &CreateEventAction{}
}
