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

type listFilesActionProps struct {
	listFoldersActionProps
}

type ListFilesAction struct{}

func (a *ListFilesAction) Name() string {
	return "List Files"
}

func (a *ListFilesAction) Description() string {
	return "Lists files in a specified directory or folder, allowing you to retrieve and process file information such as names, sizes, and timestamps."
}

func (a *ListFilesAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListFilesAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listFilesDocs,
	}
}

func (a *ListFilesAction) Icon() *string {
	return nil
}

func (a *ListFilesAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"folderId":          shared.GetFoldersInput("Folder ID", "Folder ID coming from | New Folder -> id | (or any other source)", false),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *ListFilesAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listFilesActionProps](ctx.BaseContext)
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
	q := fmt.Sprintf("%v %v", "mimeType!='application/vnd.google-apps.folder'  and ", strings.Join(qarr, " and "))

	req := driveService.Files.List().
		Fields("files(id, name, mimeType, webViewLink, kind, createdTime)").
		SupportsAllDrives(input.IncludeTeamDrives).
		Q(q)

	return req.Do()
}

func (a *ListFilesAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListFilesAction) SampleData() sdkcore.JSON {
	return []map[string]any{
		{
			"kind":     "drive#file",
			"mimeType": "image/jpeg",
			"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
			"name":     "example.jpg",
		},
	}
}

func (a *ListFilesAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListFilesAction() sdk.Action {
	return &ListFilesAction{}
}
