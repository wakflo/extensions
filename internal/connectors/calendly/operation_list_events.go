package calendly

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listEventOperationProps struct {
	User   string `json:"user"`
	Status string `json:"status"`
}

type ListEventsOperation struct {
	options *sdk.OperationInfo
}

func NewListEventsOperation() *ListEventsOperation {
	return &ListEventsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Events",
			Description: "List Events",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"user": getCurrentCalendlyUserInput("User", "Select a user", true),
				"status": autoform.NewSelectField().
					SetDisplayName("Status").
					SetDescription("Event Status").
					SetOptions(calendlyEventStatusType).
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListEventsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[listEventOperationProps](ctx)

	if input.Status == "" {
		return nil, errors.New("status is required")
	}

	if input.User == "" {
		return nil, errors.New("user is required")
	}

	reqURL := baseURL + "/scheduled_events"
	events, _ := listEvents(accessToken, reqURL, input.Status, input.User)

	return events, nil
}

func (c *ListEventsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListEventsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
