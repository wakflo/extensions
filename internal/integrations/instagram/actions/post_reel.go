package actions

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type postReelActionProps struct {
	Name string `json:"name"`
}

type PostReelAction struct{}

func (a *PostReelAction) Name() string {
	return "Post Reel"
}

func (a *PostReelAction) Description() string {
	return "Automatically posts a new reel to your social media platform, allowing you to seamlessly share engaging content with your audience."
}

func (a *PostReelAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *PostReelAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &postReelDocs,
	}
}

func (a *PostReelAction) Icon() *string {
	return nil
}

func (a *PostReelAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *PostReelAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[postReelActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *PostReelAction) Auth() *sdk.Auth {
	return nil
}

func (a *PostReelAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *PostReelAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewPostReelAction() sdk.Action {
	return &PostReelAction{}
}
