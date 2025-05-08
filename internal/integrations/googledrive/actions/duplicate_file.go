package actions

import (
	"context"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
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

// Metadata returns metadata about the action
func (a *DuplicateFileAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "duplicate_file",
		DisplayName:   "Duplicate File",
		Description:   "Duplicates one or more files and saves them with a unique identifier appended to the original file name. This action is useful when you need to create multiple copies of a file for testing, backup, or other purposes.",
		Type:          core.ActionTypeAction,
		Documentation: duplicateFileDocs,
		SampleOutput: map[string]any{
			"kind":     "drive#file",
			"mimeType": "image/jpeg",
			"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
			"name":     "example.jpg",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *DuplicateFileAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("duplicate_file", "Duplicate File")

	form.TextField("fileId", "File ID").
		Placeholder("Enter a file ID").
		Required(true).
		HelpText("The ID of the file to duplicate")

	form.TextField("fileName", "Name").
		Placeholder("Enter a file name").
		Required(true).
		HelpText("The name of the new file")

	// Add folder ID field
	shared.RegisterFoldersProp(form, "folderId", "The ID of the folder where the file will be duplicated", false)

	// Add include team drives field
	form.CheckboxField("includeTeamDrives", "Include Team Drives").
		Placeholder("Enter a value for Include Team Drives.").
		Required(false).
		HelpText("Whether to include team drives in the folder selection.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DuplicateFileAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DuplicateFileAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[duplicateFileActionProps](ctx)
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

	var parents []string
	if input.FolderID != nil {
		parents = append(parents, *input.FolderID)
	}

	in := &drive.File{
		Name:    input.FileName,
		Parents: parents,
	}

	result, err := driveService.Files.Copy(input.FileID, in).SupportsAllDrives(input.IncludeTeamDrives).Do()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewDuplicateFileAction() sdk.Action {
	return &DuplicateFileAction{}
}
