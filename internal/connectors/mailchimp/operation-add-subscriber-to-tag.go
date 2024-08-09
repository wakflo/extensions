package mailchimp

import (
	"errors"
	"log"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type addSubscriberToTagOperationProps struct {
	ListID string `json:"list-id"`
	Tags   string `json:"tag-names"`
	Email  string `json:"email"`
}

type AddSubscriberToTagOperation struct {
	options *sdk.OperationInfo
}

func NewAddSubscriberToTagOperation() sdk.IOperation {
	return &AddSubscriberToTagOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Subscriber to a tag",
			Description: "Adds a subscriber to a tag. This will fail if the user is not subscribed to the audience.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"list-id": autoform.NewShortTextField().
					SetDisplayName(" List (Audience) ID").
					SetDescription("List ID").
					SetRequired(true).
					Build(),
				"tag-names": autoform.NewLongTextField().
					SetDisplayName(" Tag Name").
					SetDescription("Tag name to add to the subscriber").
					SetRequired(true).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("Email of the subscriber").
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

func (c *AddSubscriberToTagOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[addSubscriberToTagOperationProps](ctx)
	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	tags := processTagNamesInput(input.Tags)

	err = modifySubscriberTags(accessToken, dc, input.ListID, input.Email, tags, "active")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Tag added!",
	}, nil
}

func (c *AddSubscriberToTagOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddSubscriberToTagOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
