package actions

import (
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getAllTagActionProps struct {
}

type GetAllListAction struct{}

func (a *GetAllListAction) Name() string {
	return "Get All Tag"
}

func (a *GetAllListAction) Description() string {
	return "Retrieves all tags associated with a specific entity or resource, allowing you to access and utilize tag metadata in your workflow automation processes."
}

func (a *GetAllListAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetAllListAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getAllTagDocs,
	}
}

func (a *GetAllListAction) Icon() *string {
	return nil
}

func (a *GetAllListAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (a *GetAllListAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(accessToken)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %w", err)
	}

	var result interface{}
	result, err = shared.FetchMailchimpLists(accessToken, dc)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(map[string]interface{}{
		"result": result,
	}), err
}

func (a *GetAllListAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetAllListAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"result": "Hello World!",
	}
}

func (a *GetAllListAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetAllListAction() sdk.Action {
	return &GetAllListAction{}
}
