package typeform

import (
	"errors"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type retrieveFormResponsesProps struct {
	FormID string `json:"form_id"`
}

type RetrieveFormResponsesOperation struct {
	options *sdk.OperationInfo
}

func NewRetrieveFormResponsesOperation() *RetrieveFormResponsesOperation {
	return &RetrieveFormResponsesOperation{
		options: &sdk.OperationInfo{
			Name:        "Retrieve a form response",
			Description: "Retrieve a from response from Typeform using the form ID",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"form_id": getTypeformFormsInput("Form ID", "Select a form", false),
			},
		},
	}
}

func (r *RetrieveFormResponsesOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing typeform auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[retrieveFormResponsesProps](ctx)

	if input.FormID == "" {
		return nil, errors.New("form ID is required")
	}

	resData, err := getFormResponses(accessToken, input.FormID)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

func (r *RetrieveFormResponsesOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return r.Run(ctx)
}

func (r *RetrieveFormResponsesOperation) GetInfo() *sdk.OperationInfo {
	return r.options
}
