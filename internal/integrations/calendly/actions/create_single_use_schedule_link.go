package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/calendly/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createSingleUseScheduleLinkActionProps struct {
	User          string `json:"user"`
	MaxEventCount int    `json:"max_event_count"`
	Owner         string `json:"owner"`
}

type CreateSingleUseScheduleLinkAction struct{}

func (a *CreateSingleUseScheduleLinkAction) Name() string {
	return "Create Single Use Schedule Link"
}

func (a *CreateSingleUseScheduleLinkAction) Description() string {
	return "Create a single-use schedule link that can be used to share a specific schedule with others. This link is valid only once and expires after use, ensuring secure sharing of sensitive information."
}

func (a *CreateSingleUseScheduleLinkAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateSingleUseScheduleLinkAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createSingleUseScheduleLinkDocs,
	}
}

func (a *CreateSingleUseScheduleLinkAction) Icon() *string {
	return nil
}

func (a *CreateSingleUseScheduleLinkAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"user": shared.GetCurrentCalendlyUserInput("User", "Select a user", true),
		"max_event_count": autoform.NewNumberField().
			SetDisplayName("Max Event Count").
			SetDescription("Maximum number of events that can be scheduled using the schedule link").
			SetRequired(true).
			Build(),
		"owner": shared.GetCalendlyEventTypeInput("Event Type", "Select an event type", false),
	}
}

func (a *CreateSingleUseScheduleLinkAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createSingleUseScheduleLinkActionProps](ctx.BaseContext)
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

	if input.MaxEventCount <= 0 {
		return nil, errors.New("max event count is required")
	}

	if input.Owner == "" {
		return nil, errors.New("owner is required")
	}

	scheduleLink, _ := shared.CreateSingleUseLink(accessToken, input.Owner, input.MaxEventCount)

	return scheduleLink, nil
}

func (a *CreateSingleUseScheduleLinkAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateSingleUseScheduleLinkAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateSingleUseScheduleLinkAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateSingleUseScheduleLinkAction() sdk.Action {
	return &CreateSingleUseScheduleLinkAction{}
}
