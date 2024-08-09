package mailchimp

import (
	"errors"
	"log"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type GetAllListOperation struct {
	options *sdk.OperationInfo
}

func NewGetAllListOperation() sdk.IOperation {
	return &GetAllListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get all lists ",
			Description: "Get all available lists",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetAllListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	var result interface{}
	result, err = fetchMailchimpLists(accessToken, dc)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(map[string]interface{}{
		"result": result,
	}), err
}

func (c *GetAllListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetAllListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
