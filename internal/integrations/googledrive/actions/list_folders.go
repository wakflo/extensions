package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type listFoldersActionProps struct {
	FolderID          *string `json:"folderId"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type ListFoldersAction struct{}

func (a *ListFoldersAction) Name() string {
	return "List Folders"
}

func (a *ListFoldersAction) Description() string {
	return "List Folders integration action retrieves a list of folders from a specified source, such as a cloud storage service or file system. This action allows you to access and manipulate folder structures within your workflow automation process."
}

func (a *ListFoldersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListFoldersAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listFoldersDocs,
	}
}

func (a *ListFoldersAction) Icon() *string {
	return nil
}

func (a *ListFoldersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderId":          shared.RegisterFoldersProp("Folder ID", "Folder ID coming from | New Folder -> id | (or any other source)", false),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *ListFoldersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listFoldersActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}
	var qarr []string
	if input.FolderID != nil {
		qarr = append(qarr, fmt.Sprintf("'%v' in parents", *input.FolderID))
	}
	qarr = append(qarr, "trashed = false")
	q := fmt.Sprintf("%v %v", "mimeType=='application/vnd.google-apps.folder'  and ", strings.Join(qarr, " and "))

	req := driveService.Files.List().
		Fields("files(id, name, mimeType, webViewLink, kind, createdTime)").
		SupportsAllDrives(input.IncludeTeamDrives).
		Q(q)

	return req.Do()
}

func (a *ListFoldersAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListFoldersAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListFoldersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListFoldersAction() sdk.Action {
	return &ListFoldersAction{}
}
