package actions

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

type createEventActionProps struct {
	CalendarID  string `json:"calendar_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"`
	StartTime   string `json:"start_time"`
	EndDate     string `json:"end_date"`
	EndTime     string `json:"end_time"`
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
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateEventAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_event", "Create Event")

	shared.RegisterCalendarProps(form)

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

	form.DateField("start_date", "Start Date").
		Placeholder("Event Start Date").
		HelpText("The start date of the event (YYYY-MM-DD)").
		Required(true)

	form.TextField("start_time", "Start Time").
		Placeholder("Event Start Time").
		HelpText("The start time of the event (e.g., 3pm, 10:30am, 14:00)").
		Required(true)

	form.DateField("end_date", "End Date").
		Placeholder("Event End Date").
		HelpText("The end date of the event (YYYY-MM-DD)").
		Required(true)

	form.TextField("end_time", "End Time").
		Placeholder("Event End Time").
		HelpText("The end time of the event (e.g., 5pm, 11:30am, 16:00)").
		Required(true)

	schema := form.Build()

	return schema
}

func (a *CreateEventAction) Auth() *core.AuthMetadata {
	return nil
}

// parseTime extracts hour and minute from time string
func parseTime(timeStr string) (int, int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return hour, minute, nil
}

// parseTimeInput converts various time formats to 24-hour format
func parseTimeInput(timeStr string) (string, error) {
	timeStr = strings.TrimSpace(strings.ToLower(timeStr))

	// Regular expression patterns
	// Matches: 3pm, 3:30pm, 10am, 10:45am, etc.
	ampmPattern := regexp.MustCompile(`^(\d{1,2})(?::(\d{2}))?\s*(am|pm)$`)

	// Matches: 14:00, 9:30, etc.
	hourMinPattern := regexp.MustCompile(`^(\d{1,2}):(\d{2})$`)

	// Check AM/PM format
	if matches := ampmPattern.FindStringSubmatch(timeStr); matches != nil {
		hour, _ := strconv.Atoi(matches[1])
		minute := 0
		if matches[2] != "" {
			minute, _ = strconv.Atoi(matches[2])
		}

		// Convert to 24-hour format
		if matches[3] == "pm" && hour != 12 {
			hour += 12
		} else if matches[3] == "am" && hour == 12 {
			hour = 0
		}

		return fmt.Sprintf("%02d:%02d", hour, minute), nil
	}

	// Check HH:MM format
	if matches := hourMinPattern.FindStringSubmatch(timeStr); matches != nil {
		hour, _ := strconv.Atoi(matches[1])
		minute, _ := strconv.Atoi(matches[2])
		return fmt.Sprintf("%02d:%02d", hour, minute), nil
	}

	// Check if just a number (assume hours)
	if hour, err := strconv.Atoi(timeStr); err == nil && hour >= 0 && hour <= 23 {
		return fmt.Sprintf("%02d:00", hour), nil
	}

	return "", fmt.Errorf("invalid time format: %s", timeStr)
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

	event, err := eventService.Events.Insert(input.CalendarID, &calendar.Event{
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

func NewCreateEventAction() sdk.Action {
	return &CreateEventAction{}
}
