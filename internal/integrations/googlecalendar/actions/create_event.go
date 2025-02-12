package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/googlecalendar/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *CreateEventAction) Name() string {
	return "Create Event"
}

func (a *CreateEventAction) Description() string {
	return "Create Event: Triggers the creation of a new event in your chosen calendar or scheduling system, allowing you to automate the process of setting up meetings, appointments, and other events from within your workflow."
}

func (a *CreateEventAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateEventAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createEventDocs,
	}
}

func (a *CreateEventAction) Icon() *string {
	return nil
}

func (a *CreateEventAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"calendar_id": shared.GetCalendarInput("Calendar", "select calendar", true),
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
	}
}

func (a *CreateEventAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createEventActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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
			DateTime: shared.FormatTimeString(input.Start),
		},
		End: &calendar.EventDateTime{
			DateTime: shared.FormatTimeString(input.End),
		},
	}).Do()
	return event, err
}

func (a *CreateEventAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateEventAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateEventAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateEventAction() sdk.Action {
	return &CreateEventAction{}
}
