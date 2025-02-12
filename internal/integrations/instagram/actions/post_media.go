package actions

import (
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type postMediaActionProps struct {
	Name string `json:"name"`
}

type PostMediaAction struct{}

func (a *PostMediaAction) Name() string {
	return "Post Media"
}

func (a *PostMediaAction) Description() string {
	return "Posts media content to a specified platform or service, such as social media, email, or messaging apps."
}

func (a *PostMediaAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *PostMediaAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &postMediaDocs,
	}
}

func (a *PostMediaAction) Icon() *string {
	return nil
}

func (a *PostMediaAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *PostMediaAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[postMediaActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *PostMediaAction) Auth() *sdk.Auth {
	return nil
}

func (a *PostMediaAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *PostMediaAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewPostMediaAction() sdk.Action {
	return &PostMediaAction{}
}
