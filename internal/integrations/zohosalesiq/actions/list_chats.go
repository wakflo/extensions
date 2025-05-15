package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohosalesiq/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listChatsActionProps struct {
	ScreenName string `json:"screen_name"`
}

type ListChatsAction struct{}

func (c *ListChatsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "List Chats",
		Description:   "Retrieve a list of chats from the CRM system, providing insights into customer interactions and conversations.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listChatsDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (c ListChatsAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("list_chats", "List Chats")

	form.TextField("screen_name", "Screen name").
		Placeholder("Enter a value for Screen name.").
		Required(true).
		HelpText("The screen name")

	schema := form.Build()

	return schema
}

func (c ListChatsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (c ListChatsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listChatsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	url := fmt.Sprintf("/%s/chats", input.ScreenName)

	chats, err := shared.GetZohoClient(token, url)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func NewListChatsAction() sdk.Action {
	return &ListChatsAction{}
}
