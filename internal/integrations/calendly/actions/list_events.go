package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listEventsActionProps struct {
	User   string `json:"user"`
	Status string `json:"status"`
}

type ListEventsAction struct{}

// Metadata returns metadata about the action
func (a *ListEventsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_events",
		DisplayName:   "List Events",
		Description:   "Retrieve a list of all events in your workflow, including their status and any relevant metadata. This action allows you to easily track and manage the history of your workflow's execution.",
		Type:          core.ActionTypeAction,
		Documentation: listEventsDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListEventsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_events", "List Events")

	shared.GetCurrentCalendlyUserProp("user", "User", "The user to create the schedule link for", true, form)

	form.SelectField("status", "Status").
		Placeholder("Select a status").
		Required(true).
		AddOptions(
			shared.CalendlyEventStatusType...,
		).
		HelpText("Event Status")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListEventsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListEventsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listEventsActionProps](ctx)
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

	if input.Status == "" {
		return nil, errors.New("status is required")
	}

	if input.User == "" {
		return nil, errors.New("user is required")
	}

	reqURL := shared.BaseURL + "/scheduled_events"
	events, _ := shared.ListEvents(accessToken, reqURL, input.Status, input.User)

	return events, nil
}

func NewListEventsAction() sdk.Action {
	return &ListEventsAction{}
}
