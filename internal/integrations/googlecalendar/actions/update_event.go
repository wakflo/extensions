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

type updateEventActionProps struct {
	CalendarID  string `json:"calendar_id"`
	EventID     string `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Start       string `json:"start"`
	End         string `json:"end"`
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

	form.DateTimeField("start", "start").
		Placeholder("Event Start Time").
		HelpText("The start time of the event").
		Required(false)

	form.DateTimeField("end", "end").
		Placeholder("Event end time").
		HelpText("The end time of the event").
		Required(false)

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

	if input.Start == "" {
		return nil, errors.New("start time is required")
	}

	if input.End == "" {
		return nil, errors.New("end time is required")
	}

	event, err := eventService.Events.Update(input.CalendarID, input.EventID, &calendar.Event{
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

func NewUpdateEventAction() sdk.Action {
	return &UpdateEventAction{}
}
