package actions

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type tagLocationActionProps struct {
	Name string `json:"name"`
}

type TagLocationAction struct{}

func (a *TagLocationAction) Name() string {
	return "Tag Location"
}

func (a *TagLocationAction) Description() string {
	return "Automatically assigns a specific tag to a location in your workflow, allowing you to easily identify and track locations throughout your process."
}

func (a *TagLocationAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *TagLocationAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &tagLocationDocs,
	}
}

func (a *TagLocationAction) Icon() *string {
	return nil
}

func (a *TagLocationAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *TagLocationAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[tagLocationActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *TagLocationAction) Auth() *sdk.Auth {
	return nil
}

func (a *TagLocationAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *TagLocationAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewTagLocationAction() sdk.Action {
	return &TagLocationAction{}
}
