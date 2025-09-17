package actions

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	CalendarID  string `json:"calendar_id"`
	EventID     string `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"`
	StartTime   string `json:"start_time"`
	EndDate     string `json:"end_date"`
	EndTime     string `json:"end_time"`
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
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateEventAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_event", "Update Event")

	shared.RegisterCalendarProps(form)

	shared.RegisterCalendarEventProps(form)

	form.TextField("title", "title").
		Placeholder("Event Title").
		HelpText("The title of the event.").
		Required(true)

	form.TextareaField("description", "description").
		Placeholder("Event Description").
		HelpText("The description of the event").
		Required(true)

	form.TextField("location", "location").
		Placeholder("Event Location").
		HelpText("The location of the event").
		Required(true)

	form.DateField("start_date", "start_date").
		Placeholder("Event Start Date").
		HelpText("The start date of the event (YYYY-MM-DD)").
		Required(true)

	form.TextField("start_time", "start_time").
		Placeholder("Event Start Time").
		HelpText("The start time of the event (e.g., 3pm, 10:30am, 14:00)").
		Required(true)

	form.DateField("end_date", "end_date").
		Placeholder("Event End Date").
		HelpText("The end date of the event (YYYY-MM-DD)").
		Required(true)

	form.TextField("end_time", "end_time").
		Placeholder("Event End Time").
		HelpText("The end time of the event (e.g., 5pm, 11:30am, 16:00)").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateEventAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
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

	if input.Title == "" {
		return nil, errors.New("summary is required")
	}

	if input.Description == "" {
		return nil, errors.New("description is required")
	}

	if input.Location == "" {
		return nil, errors.New("location is required")
	}

	if input.StartDate == "" {
		return nil, errors.New("start date is required")
	}

	if input.StartTime == "" {
		return nil, errors.New("start time is required")
	}

	if input.EndDate == "" {
		return nil, errors.New("end date is required")
	}

	if input.EndTime == "" {
		return nil, errors.New("end time is required")
	}

	// Parse time inputs
	startTime, err := parseTimeInput(input.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %v", err)
	}

	endTime, err := parseTimeInput(input.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %v", err)
	}

	// Extract just the date part from the input (in case it's a full datetime)
	startDateStr := input.StartDate
	if strings.Contains(startDateStr, "T") {
		startDateStr = strings.Split(startDateStr, "T")[0]
	}

	endDateStr := input.EndDate
	if strings.Contains(endDateStr, "T") {
		endDateStr = strings.Split(endDateStr, "T")[0]
	}

	// Parse the date strings
	startDateParsed, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDateParsed, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Parse the time components
	startHour, startMin, _ := parseTime(startTime)
	endHour, endMin, _ := parseTime(endTime)

	// Create full datetime with UTC timezone
	startDateTime := time.Date(
		startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(),
		startHour, startMin, 0, 0, time.UTC,
	).Format(time.RFC3339)

	endDateTime := time.Date(
		endDateParsed.Year(), endDateParsed.Month(), endDateParsed.Day(),
		endHour, endMin, 0, 0, time.UTC,
	).Format(time.RFC3339)

	event, err := eventService.Events.Update(input.CalendarID, input.EventID, &calendar.Event{
		Summary:     input.Title,
		Description: input.Description,
		Location:    input.Location,
		Start: &calendar.EventDateTime{
			DateTime: startDateTime,
		},
		End: &calendar.EventDateTime{
			DateTime: endDateTime,
		},
	}).Do()
	return event, err
}

func NewUpdateEventAction() sdk.Action {
	return &UpdateEventAction{}
}
