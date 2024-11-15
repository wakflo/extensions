package zohosalesiq

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listChatsOperationProps struct {
	ScreenName string `json:"screen_name"`
}

type ListChatsOperation struct {
	options *sdk.OperationInfo
}

func NewListChatsOperation() sdk.IOperation {
	return &ListChatsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Chats",
			Description: "Get list of chats",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"screen_name": autoform.NewShortTextField().
					SetDisplayName("Screen name").
					SetDescription("Screen name").
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

func (c *ListChatsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[listChatsOperationProps](ctx)

	url := fmt.Sprintf("/%s/chats", input.ScreenName)

	chats, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (c *ListChatsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListChatsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
