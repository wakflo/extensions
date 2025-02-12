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

type uploadFileActionProps struct {
	FileName          string  `json:"fileName"`
	File              string  `json:"file"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type UploadFileAction struct{}

func (a *UploadFileAction) Name() string {
	return "Upload File"
}

func (a *UploadFileAction) Description() string {
	return "Upload File: This integration action allows you to upload files from various sources such as cloud storage services, local file systems, or email attachments to your workflow. You can specify the file type, size limit, and other parameters to control the upload process. The uploaded file is then stored in a designated location within your workflow, making it easily accessible for further processing or analysis."
}

func (a *UploadFileAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UploadFileAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &uploadFileDocs,
	}
}

func (a *UploadFileAction) Icon() *string {
	return nil
}

func (a *UploadFileAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"fileName": autoform.NewShortTextField().
			SetDisplayName("Name").
			SetDescription("The name of the new file").
			SetRequired(true).
			Build(),
		"file": autoform.NewFileField().
			SetDisplayName("File").
			SetDescription("The file URL or base64 to upload").
			SetRequired(true).
			Build(),
		"parentFolder":      shared.GetParentFoldersInput(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *UploadFileAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[uploadFileActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	fileData, err := sdk.StringToFile(input.File)
	if err != nil {
		return nil, err
	}

	var parents []string
	if input.ParentFolder != nil {
		parents = append(parents, *input.ParentFolder)
	}

	in := &drive.File{
		MimeType: fileData.Mime,
		Name:     input.FileName,
		Parents:  parents,
	}

	return driveService.Files.Create(in).
		Media(fileData.Data).
		SupportsAllDrives(input.IncludeTeamDrives).
		Do()
}

func (a *UploadFileAction) Auth() *sdk.Auth {
	return nil
}

func (a *UploadFileAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UploadFileAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUploadFileAction() sdk.Action {
	return &UploadFileAction{}
}
