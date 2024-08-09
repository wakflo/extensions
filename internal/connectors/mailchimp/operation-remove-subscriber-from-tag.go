package mailchimp

import (
	"errors"
	"log"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type removeSubscriberToTagOperationProps struct {
	ListID string `json:"list-id"`
	Tags   string `json:"tag-names"`
	Email  string `json:"email"`
}

type RemoveSubscriberToTagOperation struct {
	options *sdk.OperationInfo
}

func NewRemoveSubscriberToTagOperation() sdk.IOperation {
	return &RemoveSubscriberToTagOperation{
		options: &sdk.OperationInfo{
			Name:        "Remove a Subscriber from a tag",
			Description: "Removes a subscriber from a tag.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"list-id": autoform.NewShortTextField().
					SetDisplayName(" List (Audience) ID").
					SetDescription("List ID").
					SetRequired(true).
					Build(),
				"tag-names": autoform.NewLongTextField().
					SetDisplayName(" Tag name").
					SetDescription("Tag name to remove from the subscriber").
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

func (c *RemoveSubscriberToTagOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[removeSubscriberToTagOperationProps](ctx)
	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	tags := processTagNamesInput(input.Tags)

	err = modifySubscriberTags(accessToken, dc, input.ListID, input.Email, tags, "inactive")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Tag Removed!",
	}, nil
}

func (c *RemoveSubscriberToTagOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *RemoveSubscriberToTagOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
