package actions

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type uploadMediaFromUrlActionProps struct {
	Name string `json:"name"`
}

type UploadMediaFromUrlAction struct{}

func (a *UploadMediaFromUrlAction) Name() string {
	return "Upload Media From URL"
}

func (a *UploadMediaFromUrlAction) Description() string {
	return "Uploads media files from a specified URL to your workflow's storage location, allowing you to easily incorporate external content into your automated processes. This integration action is particularly useful when working with cloud-based file sharing services or public media repositories. Simply provide the URL of the media file you wish to upload, and our software will handle the rest, saving the file to a designated location for use in subsequent workflow steps."
}

func (a *UploadMediaFromUrlAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UploadMediaFromUrlAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &uploadMediaFromUrlDocs,
	}
}

func (a *UploadMediaFromUrlAction) Icon() *string {
	return nil
}

func (a *UploadMediaFromUrlAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *UploadMediaFromUrlAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[uploadMediaFromUrlActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *UploadMediaFromUrlAction) Auth() *sdk.Auth {
	return nil
}

func (a *UploadMediaFromUrlAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UploadMediaFromUrlAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUploadMediaFromUrlAction() sdk.Action {
	return &UploadMediaFromUrlAction{}
}
