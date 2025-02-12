package actions

import (
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/zohosalesiq/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listChatsActionProps struct {
	ScreenName string `json:"screen_name"`
}

type ListChatsAction struct{}

func (c *ListChatsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c ListChatsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c ListChatsAction) Name() string {
	return "List Chats"
}

func (c ListChatsAction) Description() string {
	return "Get list of chats"
}

func (c ListChatsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listChatsDocs,
	}
}

func (c ListChatsAction) Icon() *string {
	return nil
}

func (c ListChatsAction) SampleData() sdkcore.JSON {
	return nil
}

func (c ListChatsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"screen_name": autoform.NewShortTextField().
			SetDisplayName("Screen name").
			SetDescription("Screen name").
			SetRequired(true).
			Build(),
	}
}

func (c ListChatsAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c ListChatsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listChatsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/%s/chats", input.ScreenName)

	chats, err := shared.GetZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func NewListChatsAction() sdk.Action {
	return &ListChatsAction{}
}
