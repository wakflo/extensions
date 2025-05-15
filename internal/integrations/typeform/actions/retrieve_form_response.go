package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/typeform/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type retrieveFormResponseActionProps struct {
	FormID string `json:"form_id"`
}

type RetrieveFormResponseAction struct{}

func (a *RetrieveFormResponseAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Retrieve a form response",
		Description:   "Retrieve a from response from Typeform using the form ID",
		Type:          core.ActionTypeAction,
		Documentation: retrieveFormResponseDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
		Icon:          "mingcute:report-forms-fill",
	}
}

func (a *RetrieveFormResponseAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("retrieve_form_response", "Retrieve a form response")
	shared.RegisterTypeformFormsProps(form, "Form ID", "Select a form", true)

	schema := form.Build()
	return schema
}

func (a *RetrieveFormResponseAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[retrieveFormResponseActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if input.FormID == "" {
		return nil, errors.New("form ID is required")
	}
	token := *&authCtx.Token.AccessToken

	resData, err := shared.GetFormResponses(token, input.FormID)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

func (a *RetrieveFormResponseAction) Auth() *core.AuthMetadata {
	return nil
}

func NewRetrieveFormResponseAction() sdk.Action {
	return &RetrieveFormResponseAction{}
}
