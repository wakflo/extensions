package actions

import (
	"errors"

	"github.com/wakflo/extensions/internal/integrations/typeform/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type retrieveFormResponseActionProps struct {
	FormID string `json:"form_id"`
}

type RetrieveFormResponseAction struct{}

func (a *RetrieveFormResponseAction) Name() string {
	return "Retrieve a form response"
}

func (a *RetrieveFormResponseAction) Description() string {
	return "Retrieve a from response from Typeform using the form ID"
}

func (a *RetrieveFormResponseAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *RetrieveFormResponseAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &retrieveFormResponseDocs,
	}
}

func (a *RetrieveFormResponseAction) Icon() *string {
	icon := "mingcute:report-forms-fill"
	return &icon
}

func (a *RetrieveFormResponseAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"form_id": shared.GetTypeformFormsInput("Form ID", "Select a form", true),
	}
}

func (a *RetrieveFormResponseAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[retrieveFormResponseActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	if input.FormID == "" {
		return nil, errors.New("form ID is required")
	}

	resData, err := shared.GetFormResponses(accessToken, input.FormID)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

func (a *RetrieveFormResponseAction) Auth() *sdk.Auth {
	return nil
}

func (a *RetrieveFormResponseAction) SampleData() sdkcore.JSON {
	return nil
}

func (a *RetrieveFormResponseAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewRetrieveFormResponseAction() sdk.Action {
	return &RetrieveFormResponseAction{}
}
