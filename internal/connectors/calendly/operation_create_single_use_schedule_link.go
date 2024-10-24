package calendly

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createSingleUseScheduleLinkOperationProps struct {
	User          string `json:"user"`
	MaxEventCount int    `json:"max_event_count"`
	Owner         string `json:"owner"`
}

type CreateSingleUseScheduleLinkOperation struct {
	options *sdk.OperationInfo
}

func NewCreateSingleUseScheduleLinkOperation() *CreateSingleUseScheduleLinkOperation {
	return &CreateSingleUseScheduleLinkOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Single Use Schedule Link",
			Description: "Create a single use schedule link",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"user": getCurrentCalendlyUserInput("User", "Select a user", true),
				"max_event_count": autoform.NewNumberField().
					SetDisplayName("Max Event Count").
					SetDescription("Maximum number of events that can be scheduled using the schedule link").
					SetRequired(true).
					Build(),
				"owner": getCalendlyEventTypeInput("Event Type", "Select an event type", false),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateSingleUseScheduleLinkOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[createSingleUseScheduleLinkOperationProps](ctx)

	if input.User == "" {
		return nil, errors.New("user is required")
	}

	if input.MaxEventCount <= 0 {
		return nil, errors.New("max event count is required")
	}

	if input.Owner == "" {
		return nil, errors.New("owner is required")
	}

	scheduleLink, _ := createSingleUseLink(accessToken, input.Owner, input.MaxEventCount)

	return scheduleLink, nil
}

func (c *CreateSingleUseScheduleLinkOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateSingleUseScheduleLinkOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
