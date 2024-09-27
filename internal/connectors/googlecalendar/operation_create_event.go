package googlecalendar

import (
	"context"
	"errors"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createEventOperationProps struct {
	CalendarID  string `json:"calendar_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Start       string `json:"start"`
	End         string `json:"end"`
}

type CreateEventOperation struct {
	options *sdk.OperationInfo
}

func NewCreateEventOperation() *CreateEventOperation {
	return &CreateEventOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Event",
			Description: "create a calendar event",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"calendar_id": getCalendarInput("Calendar", "select calendar", true),
				"title": autoform.NewShortTextField().
					SetDisplayName("Event Title").
					SetDescription("The title of the event.").
					SetRequired(true).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Event Description").
					SetDescription("The description of the event").
					SetRequired(true).
					Build(),
				"location": autoform.NewShortTextField().
					SetDisplayName("Event Location").
					SetDescription("The location of the event").
					SetRequired(true).
					Build(),
				"start": autoform.NewDateTimeField().
					SetDisplayName("Event Start Time").
					SetDescription("The start time of the event").
					SetRequired(true).
					Build(),
				"end": autoform.NewDateTimeField().
					SetDisplayName("Event end time").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateEventOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[createEventOperationProps](ctx)
	eventService, err := calendar.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))

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

	if input.Start == "" {
		return nil, errors.New("start time is required")
	}

	if input.End == "" {
		return nil, errors.New("end time is required")
	}

	event, err := eventService.Events.Insert(input.CalendarID, &calendar.Event{
		Summary:     input.Title,
		Description: input.Description,
		Location:    input.Location,
		Start: &calendar.EventDateTime{
			DateTime: formatTimeString(input.Start),
		},
		End: &calendar.EventDateTime{
			DateTime: formatTimeString(input.End),
		},
	}).Do()
	return event, err
}

func (c *CreateEventOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateEventOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
