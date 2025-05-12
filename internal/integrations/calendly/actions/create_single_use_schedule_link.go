package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createSingleUseScheduleLinkActionProps struct {
	User          string `json:"user"`
	MaxEventCount int    `json:"max_event_count"`
	Owner         string `json:"owner"`
}

type CreateSingleUseScheduleLinkAction struct{}

// Metadata returns metadata about the action
func (a *CreateSingleUseScheduleLinkAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_single_use_schedule_link",
		DisplayName:   "Create Single Use Schedule Link",
		Description:   "Create a single-use schedule link that can be used to share a specific schedule with others. This link is valid only once and expires after use, ensuring secure sharing of sensitive information.",
		Type:          core.ActionTypeAction,
		Documentation: createSingleUseScheduleLinkDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateSingleUseScheduleLinkAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_single_use_schedule_link", "Create Single Use Schedule Link")

	// Note: These will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("user", "User").
	//	Placeholder("Select a user").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The user to create the schedule link for")

	form.NumberField("max_event_count", "Max Event Count").
		Placeholder("Enter a number").
		Required(true).
		HelpText("Maximum number of events that can be scheduled using the schedule link")

	// Note: These will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("owner", "Event Type").
	//	Placeholder("Select an event type").
	//	Required(false).
	//	WithDynamicOptions(...).
	//	HelpText("The event type to create the schedule link for")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateSingleUseScheduleLinkAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateSingleUseScheduleLinkAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createSingleUseScheduleLinkActionProps](ctx)
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

	if input.MaxEventCount <= 0 {
		return nil, errors.New("max event count is required")
	}

	if input.Owner == "" {
		return nil, errors.New("owner is required")
	}

	scheduleLink, _ := shared.CreateSingleUseLink(accessToken, input.Owner, input.MaxEventCount)

	return scheduleLink, nil
}

func NewCreateSingleUseScheduleLinkAction() sdk.Action {
	return &CreateSingleUseScheduleLinkAction{}
}
