package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getAllTagActionProps struct{}

type GetAllListAction struct{}

func (a *GetAllListAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Get All List",
		Description:   "Retrieves all lists associated with a specific entity or resource, allowing you to access and utilize list metadata in your workflow automation processes.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getAllTagDocs,
		SampleOutput: map[string]any{
			"success": "List retrieved!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetAllListAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_all_list", "Get All List")

	schema := form.Build()

	return schema
}

func (a *GetAllListAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(token)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %w", err)
	}

	var result interface{}
	result, err = shared.FetchMailchimpLists(token, dc)
	if err != nil {
		return nil, err
	}

	return sdkcore.JSON(map[string]interface{}{
		"result": result,
	}), err
}

func (a *GetAllListAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetAllListAction() sdk.Action {
	return &GetAllListAction{}
}
