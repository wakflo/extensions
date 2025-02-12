package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listEventsActionProps struct {
	User   string `json:"user"`
	Status string `json:"status"`
}

type ListEventsAction struct{}

func (a *ListEventsAction) Name() string {
	return "List Events"
}

func (a *ListEventsAction) Description() string {
	return "Retrieve a list of all events in your workflow, including their status and any relevant metadata. This action allows you to easily track and manage the history of your workflow's execution."
}

func (a *ListEventsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListEventsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listEventsDocs,
	}
}

func (a *ListEventsAction) Icon() *string {
	return nil
}

func (a *ListEventsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"user": shared.GetCurrentCalendlyUserInput("User", "Select a user", true),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetDescription("Event Status").
			SetOptions(shared.CalendlyEventStatusType).
			SetRequired(true).
			Build(),
	}
}

func (a *ListEventsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listEventsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

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

func (a *ListEventsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListEventsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListEventsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListEventsAction() sdk.Action {
	return &ListEventsAction{}
}
