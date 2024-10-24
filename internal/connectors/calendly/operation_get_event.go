package calendly

import (
	"errors"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getEventOperatioinProps struct {
	User    string `json:"user"`
	EventID string `json:"event-id"`
}

type GetEventOperation struct {
	options *sdk.OperationInfo
}

func NewGetEventOperation() *GetEventOperation {
	return &GetEventOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Event",
			Description: "Get Event",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"user":     getCurrentCalendlyUserInput("User", "Select a user", true),
				"event-id": getCalendlyEventInput("Event", "Select event", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetEventOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[getEventOperatioinProps](ctx)

	println("event-ud", input.EventID)
	println("accessToken", accessToken)

	if input.User == "" {
		return nil, errors.New("user is required")
	}

	if input.EventID == "" {
		return nil, errors.New("event id is required")
	}

	event, _ := getEvent(accessToken, input.EventID)

	return event, nil
}

func (c *GetEventOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetEventOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
