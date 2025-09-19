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

type updateEventActionProps struct {
	CalendarID    string `json:"calendar_id"`
	EventID       string `json:"event_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Location      string `json:"location"`
	StartDateTime string `json:"start_datetime"`
	EndDateTime   string `json:"end_datetime"`
}

type UpdateEventAction struct{}

// Metadata returns metadata about the action
func (a *UpdateEventAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_event",
		DisplayName:   "Update Event",
		Description:   "Updates an existing event in your workflow, allowing you to modify or refresh information as needed. This action enables real-time updates and ensures that all connected workflows and integrations reflect the latest changes.",
		Type:          core.ActionTypeAction,
		Documentation: updateEventDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"id":      "event_123",
			"summary": "Updated Meeting",
			"start": map[string]any{
				"dateTime": "2025-09-23T10:00:00Z",
			},
			"end": map[string]any{
				"dateTime": "2025-09-23T11:00:00Z",
			},
			"htmlLink": "https://calendar.google.com/event?eid=...",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateEventAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_event", "Update Event")

	shared.RegisterCalendarProps(form)
	shared.RegisterCalendarEventProps(form)

	form.TextField("title", "Title").
		Placeholder("Event Title").
		HelpText("The title of the event").
		Required(false)

	form.TextareaField("description", "Description").
		Placeholder("Event Description").
		HelpText("The description of the event").
		Required(false)

	form.TextField("location", "Location").
		Placeholder("Event Location").
		HelpText("The location of the event").
		Required(false)

	form.DateField("start_datetime", "Start Date & Time").
		Placeholder("Select start date and time").
		HelpText("The start date and time of the event").
		Required(false)

	form.DateField("end_datetime", "End Date & Time").
		Placeholder("Select end date and time").
		HelpText("The end date and time of the event").
		Required(false)

	schema := form.Build()

	return schema
}

func (a *UpdateEventAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *UpdateEventAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateEventActionProps](ctx)
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
		return nil, errors.New("calendar is required")
	}

	if input.EventID == "" {
		return nil, errors.New("event id is required")
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

	startDateTimeStr := startDateTime.Format(time.RFC3339)
	endDateTimeStr := endDateTime.Format(time.RFC3339)

	event, err := eventService.Events.Update(input.CalendarID, input.EventID, &calendar.Event{
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
		return nil, fmt.Errorf("failed to update event: %v", err)
	}

	return event, nil
}

func NewUpdateEventAction() sdk.Action {
	return &UpdateEventAction{}
}
