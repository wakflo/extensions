package actions

import (
	"context"

	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type readFileContentActionProps struct {
	FileID   string  `json:"fileId"`
	FileName *string `json:"fileName"`
}

type ReadFileContentAction struct{}

func (a *ReadFileContentAction) Name() string {
	return "Read File Content"
}

func (a *ReadFileContentAction) Description() string {
	return "Reads the content of a specified file and returns it as a string or binary data, depending on the file type. This action is useful when you need to extract information from a file or process its contents in your workflow automation."
}

func (a *ReadFileContentAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ReadFileContentAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &readFileContentDocs,
	}
}

func (a *ReadFileContentAction) Icon() *string {
	return nil
}

func (a *ReadFileContentAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"fileId": autoform.NewShortTextField().
			SetDisplayName("File ID").
			SetDescription("File ID coming from | New File -> id |").
			SetRequired(true).
			Build(),
		"fileName": autoform.NewShortTextField().
			SetDisplayName("File Name").
			SetDescription("Destination File name").
			SetRequired(true).
			Build(),
	}
}

func (a *ReadFileContentAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[readFileContentActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	return shared.DownloadFile(&ctx.BaseContext, driveService, input.FileID, input.FileName)
}

func (a *ReadFileContentAction) Auth() *sdk.Auth {
	return nil
}

func (a *ReadFileContentAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ReadFileContentAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewReadFileContentAction() sdk.Action {
	return &ReadFileContentAction{}
}
