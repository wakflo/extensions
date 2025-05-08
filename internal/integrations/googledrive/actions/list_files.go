package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type listFilesActionProps struct {
	listFoldersActionProps
}

type ListFilesAction struct{}

// Metadata returns metadata about the action
func (a *ListFilesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_files",
		DisplayName:   "List Files",
		Description:   "Lists files in a specified directory or folder, allowing you to retrieve and process file information such as names, sizes, and timestamps.",
		Type:          core.ActionTypeAction,
		Documentation: listFilesDocs,
		SampleOutput: map[string]any{
			"files": []map[string]any{
				{
					"kind":     "drive#file",
					"mimeType": "image/jpeg",
					"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
					"name":     "example.jpg",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListFilesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_files", "List Files")

	// Add folder ID field
	shared.RegisterFoldersProp(form, "folderId", "Folder ID coming from | New Folder -> id | (or any other source)", false)

	// Add include team drives field
	form.CheckboxField("includeTeamDrives", "Include Team Drives").
		Placeholder("Enter a value for Include Team Drives.").
		Required(false).
		HelpText("Whether to include team drives in the folder selection.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListFilesAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListFilesAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listFilesActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
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

	result, err := req.Do()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewListFilesAction() sdk.Action {
	return &ListFilesAction{}
}
