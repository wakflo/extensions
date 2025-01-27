package actions

import (
	"context"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type createFolderActionProps struct {
	FolderName        string  `json:"folderName"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type CreateFolderAction struct {
}

func (c CreateFolderAction) Name() string {
	return "Create Folder"
}

func (c CreateFolderAction) Description() string {
	return "Create a new folder in Google Drive"
}

func (c CreateFolderAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &createFolderDocs,
	}
}

func (c CreateFolderAction) Icon() *string {
	return nil
}

func (c CreateFolderAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c CreateFolderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderName": autoform.NewShortTextField().
			SetDisplayName("Folder name").
			SetDescription("The name of the new folder.").
			SetRequired(true).
			Build(),
		"parentFolder":      shared.GetParentFoldersInput(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (c CreateFolderAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c CreateFolderAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[createFolderActionProps](ctx.BaseContext)
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
		MimeType: "application/vnd.google-apps.folder",
		Name:     input.FolderName,
		Parents:  parents,
	}).
		Fields("id, name, mimeType, webViewLink, kind, createdTime").
		SupportsAllDrives(input.IncludeTeamDrives).
		Do()
	return folder, err
}

func NewCreateFolderAction() integration.Action {
	return &CreateFolderAction{}
}
