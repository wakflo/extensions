package actions

import (
	"context"
	"errors"

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
	Start       string `json:"start"`
	End         string `json:"end"`
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

	form.DateTimeField("start", "start").
		Placeholder("Event Start Time").
		HelpText("The start time of the event").
		Required(true)

	form.DateTimeField("end", "end").
		Placeholder("Event end time").
		HelpText("The end time of the event").
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
			DateTime: shared.FormatTimeString(input.Start),
		},
		End: &calendar.EventDateTime{
			DateTime: shared.FormatTimeString(input.End),
		},
	}).Do()
	return event, err
}

func NewCreateEventAction() sdk.Action {
	return &CreateEventAction{}
}
