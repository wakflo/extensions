package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getEventActionProps struct {
	User    string `json:"user"`
	EventID string `json:"event-id"`
}

type GetEventAction struct{}

func (a *GetEventAction) Name() string {
	return "Get Event"
}

func (a *GetEventAction) Description() string {
	return "Retrieves an event from the system, allowing you to access and utilize event data in your workflow automation."
}

func (a *GetEventAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetEventAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getEventDocs,
	}
}

func (a *GetEventAction) Icon() *string {
	return nil
}

func (a *GetEventAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"user":     shared.GetCurrentCalendlyUserInput("User", "Select a user", true),
		"event-id": shared.GetCalendlyEventInput("Event", "Select event", true),
	}
}

func (a *GetEventAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getEventActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	if input.User == "" {
		return nil, errors.New("user is required")
	}

	if input.EventID == "" {
		return nil, errors.New("event id is required")
	}

	event, _ := shared.GetEvent(accessToken, input.EventID)

	return event, nil
}

func (a *GetEventAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetEventAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetEventAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetEventAction() sdk.Action {
	return &GetEventAction{}
}
