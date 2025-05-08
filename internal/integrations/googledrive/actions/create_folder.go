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

type createFolderActionProps struct {
	FolderName        string  `json:"folderName"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type CreateFolderAction struct{}

func (a *CreateFolderAction) Name() string {
	return "Create Folder"
}

func (a *CreateFolderAction) Description() string {
	return "Creates a new folder in the specified location, allowing you to organize and structure your files and data within your workflow."
}

func (a *CreateFolderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateFolderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createFolderDocs,
	}
}

func (a *CreateFolderAction) Icon() *string {
	return nil
}

func (a *CreateFolderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderName": autoform.NewShortTextField().
			SetDisplayName("Folder name").
			SetDescription("The name of the new folder.").
			SetRequired(true).
			Build(),
		"parentFolder":      shared.RegisterParentFoldersProp(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *CreateFolderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createFolderActionProps](ctx.BaseContext)
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

func (a *CreateFolderAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateFolderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"kind":     "drive#file",
		"mimeType": "image/jpeg",
		"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
		"name":     "example.jpg",
	}
}

func (a *CreateFolderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateFolderAction() sdk.Action {
	return &CreateFolderAction{}
}
