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

func (a *UpdateEventAction) Name() string {
	return "Update Event"
}

func (a *UpdateEventAction) Description() string {
	return "Updates an existing event in your workflow, allowing you to modify or refresh information as needed. This action enables real-time updates and ensures that all connected workflows and integrations reflect the latest changes."
}

func (a *UpdateEventAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateEventAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateEventDocs,
	}
}

func (a *UpdateEventAction) Icon() *string {
	return nil
}

func (a *UpdateEventAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"calendar_id": shared.GetCalendarInput("Calendar", "select calendar", true),
		"event_id":    shared.GetCalendarEventIDInput("Event Id", "select event", true),
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
			SetRequired(false).
			Build(),
		"end": autoform.NewDateTimeField().
			SetDisplayName("Event end time").
			SetRequired(false).
			Build(),
	}
}

func (a *UpdateEventAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateEventActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	eventService, err := calendar.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
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

func (a *UpdateEventAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateEventAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateEventAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateEventAction() sdk.Action {
	return &UpdateEventAction{}
}
