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

type duplicateFileActionProps struct {
	FileID            string  `json:"fileId"`
	FileName          string  `json:"fileName"`
	FolderID          *string `json:"folderId"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type DuplicateFileAction struct{}

func (a *DuplicateFileAction) Name() string {
	return "Duplicate File"
}

func (a *DuplicateFileAction) Description() string {
	return "Duplicates one or more files and saves them with a unique identifier appended to the original file name. This action is useful when you need to create multiple copies of a file for testing, backup, or other purposes."
}

func (a *DuplicateFileAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *DuplicateFileAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &duplicateFileDocs,
	}
}

func (a *DuplicateFileAction) Icon() *string {
	return nil
}

func (a *DuplicateFileAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"fileId": autoform.NewShortTextField().
			SetDisplayName("File ID").
			SetDescription("The ID of the file to duplicate").
			SetRequired(true).
			Build(),
		"fileName": autoform.NewShortTextField().
			SetDisplayName("Name").
			SetDescription("The name of the new file").
			SetRequired(true).
			Build(),
		"folderId":          shared.GetFoldersInput("Folder ID", "The ID of the folder where the file will be duplicated", false),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *DuplicateFileAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[duplicateFileActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var parents []string
	if input.FolderID != nil {
		parents = append(parents, *input.FolderID)
	}

	in := &drive.File{
		Name:    input.FileName,
		Parents: parents,
	}

	return driveService.Files.Copy(input.FileID, in).SupportsAllDrives(input.IncludeTeamDrives).Do()
}

func (a *DuplicateFileAction) Auth() *sdk.Auth {
	return nil
}

func (a *DuplicateFileAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"kind":     "drive#file",
		"mimeType": "image/jpeg",
		"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
		"name":     "example.jpg",
	}
}

func (a *DuplicateFileAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDuplicateFileAction() sdk.Action {
	return &DuplicateFileAction{}
}
