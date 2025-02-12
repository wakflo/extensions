package actions

import (
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type postMediaFromUrlActionProps struct {
	Name string `json:"name"`
}

type PostMediaFromUrlAction struct{}

func (a *PostMediaFromUrlAction) Name() string {
	return "Post Media From URL"
}

func (a *PostMediaFromUrlAction) Description() string {
	return "Posts media from a specified URL to your workflow's media library, allowing you to easily incorporate external content into your automated processes."
}

func (a *PostMediaFromUrlAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *PostMediaFromUrlAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &postMediaFromUrlDocs,
	}
}

func (a *PostMediaFromUrlAction) Icon() *string {
	return nil
}

func (a *PostMediaFromUrlAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *PostMediaFromUrlAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[postMediaFromUrlActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *PostMediaFromUrlAction) Auth() *sdk.Auth {
	return nil
}

func (a *PostMediaFromUrlAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *PostMediaFromUrlAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewPostMediaFromUrlAction() sdk.Action {
	return &PostMediaFromUrlAction{}
}
