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

type createFileActionProps struct {
	FileName          string  `json:"fileName"`
	Content           []byte  `json:"content"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type CreateFileAction struct{}

func (a *CreateFileAction) Name() string {
	return "Create File"
}

func (a *CreateFileAction) Description() string {
	return "Creates a new file with a specified name and content, allowing you to store and manage data within your workflow."
}

func (a *CreateFileAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateFileAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createFileDocs,
	}
}

func (a *CreateFileAction) Icon() *string {
	return nil
}

func (a *CreateFileAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("File name").
			SetDescription("The name of the new file with extension.").
			SetRequired(true).
			Build(),
		"content": autoform.NewFileField().
			SetDisplayName("File").
			SetDescription("file to upload.").
			SetRequired(true).
			Build(),
		"parentFolder":      shared.GetParentFoldersInput(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *CreateFileAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createFileActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var parents []string
	if input.ParentFolder != nil {
		parents = append(parents, *input.ParentFolder)
	}

	folder, err := driveService.Files.Create(&drive.File{
		MimeType: "application/vnd.google-apps.file",
		Name:     input.FileName,
		Parents:  parents,
	}).
		Fields("id, name, mimeType, webViewLink, kind, createdTime").
		SupportsAllDrives(input.IncludeTeamDrives).
		Do()
	return folder, err
}

func (a *CreateFileAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateFileAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"kind":     "drive#file",
		"mimeType": "image/jpeg",
		"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
		"name":     "example.jpg",
	}
}

func (a *CreateFileAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateFileAction() sdk.Action {
	return &CreateFileAction{}
}
