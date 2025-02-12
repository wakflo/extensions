package actions

import (
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type postSingleMediaActionProps struct {
	Name string `json:"name"`
}

type PostSingleMediaAction struct{}

func (a *PostSingleMediaAction) Name() string {
	return "Post Single Media"
}

func (a *PostSingleMediaAction) Description() string {
	return "Posts a single media item (image, video, or audio) to a specified platform or service, such as social media, email, or messaging apps."
}

func (a *PostSingleMediaAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *PostSingleMediaAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &postSingleMediaDocs,
	}
}

func (a *PostSingleMediaAction) Icon() *string {
	return nil
}

func (a *PostSingleMediaAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *PostSingleMediaAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[postSingleMediaActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *PostSingleMediaAction) Auth() *sdk.Auth {
	return nil
}

func (a *PostSingleMediaAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *PostSingleMediaAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewPostSingleMediaAction() sdk.Action {
	return &PostSingleMediaAction{}
}
