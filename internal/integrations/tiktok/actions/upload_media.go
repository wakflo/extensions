package actions

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type uploadMediaActionProps struct {
	Name string `json:"name"`
}

type UploadMediaAction struct{}

func (a *UploadMediaAction) Name() string {
	return "Upload Media"
}

func (a *UploadMediaAction) Description() string {
	return "Upload Media: Automatically upload media files (images, videos, documents) from your computer or cloud storage to your desired destination, such as a database, file system, or content management system. This integration action streamlines the process of sharing and storing multimedia assets, reducing manual effort and increasing collaboration efficiency."
}

func (a *UploadMediaAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UploadMediaAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &uploadMediaDocs,
	}
}

func (a *UploadMediaAction) Icon() *string {
	return nil
}

func (a *UploadMediaAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *UploadMediaAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[uploadMediaActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *UploadMediaAction) Auth() *sdk.Auth {
	return nil
}

func (a *UploadMediaAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UploadMediaAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUploadMediaAction() sdk.Action {
	return &UploadMediaAction{}
}
