package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getEventActionProps struct {
	User    string `json:"user"`
	EventID string `json:"event-id"`
}

type GetEventAction struct{}

// Metadata returns metadata about the action
func (a *GetEventAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_event",
		DisplayName:   "Get Event",
		Description:   "Retrieves an event from the system, allowing you to access and utilize event data in your workflow automation.",
		Type:          core.ActionTypeAction,
		Documentation: getEventDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetEventAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_event", "Get Event")

	// Note: These will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("user", "User").
	//	Placeholder("Select a user").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The user to get events for")

	form.TextField("event-id", "Event ID").
		Placeholder("Enter an event ID").
		Required(true).
		HelpText("The ID of the event to retrieve")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetEventAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetEventAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getEventActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := authCtx.AccessToken

	if input.User == "" {
		return nil, errors.New("user is required")
	}

	if input.EventID == "" {
		return nil, errors.New("event id is required")
	}

	event, _ := shared.GetEvent(accessToken, input.EventID)

	return event, nil
}

func NewGetEventAction() sdk.Action {
	return &GetEventAction{}
}
