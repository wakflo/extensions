package actions

import (
	"context"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type getFileActionProps struct {
	FileID string `json:"fileId"`
}

type GetFileAction struct{}

func (a *GetFileAction) Name() string {
	return "Get File"
}

func (a *GetFileAction) Description() string {
	return "Retrieves a file from a specified location and makes it available for further processing in the workflow."
}

func (a *GetFileAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetFileAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getFileDocs,
	}
}

func (a *GetFileAction) Icon() *string {
	return nil
}

func (a *GetFileAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"fileId": autoform.NewShortTextField().
			SetDisplayName("File / Folder ID").
			SetDescription("The ID of the file/folder to search for.").
			SetRequired(true).
			Build(),
	}
}

func (a *GetFileAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getFileActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	file, err := driveService.Files.Get(input.FileID).
		Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()
	return file, err
}

func (a *GetFileAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetFileAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"kind":     "drive#file",
		"mimeType": "image/jpeg",
		"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
		"name":     "example.jpg",
	}
}

func (a *GetFileAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetFileAction() sdk.Action {
	return &GetFileAction{}
}
