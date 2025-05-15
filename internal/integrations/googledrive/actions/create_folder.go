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

type createFolderActionProps struct {
	FolderName        string  `json:"folderName"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type CreateFolderAction struct{}

// Metadata returns metadata about the action
func (a *CreateFolderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_folder",
		DisplayName:   "Create Folder",
		Description:   "Creates a new folder in the specified location, allowing you to organize and structure your files and data within your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: createFolderDocs,
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
func (a *CreateFolderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_folder", "Create Folder")

	form.TextField("folderName", "Folder name").
		Placeholder("Enter a folder name").
		Required(true).
		HelpText("The name of the new folder.")

	shared.RegisterParentFoldersProp(form)

	// Add include team drives field
	form.CheckboxField("includeTeamDrives", "Include Team Drives").
		Placeholder("Enter a value for Include Team Drives.").
		Required(false).
		HelpText("Whether to include team drives in the folder selection.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateFolderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateFolderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createFolderActionProps](ctx)
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
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func NewCreateFolderAction() sdk.Action {
	return &CreateFolderAction{}
}
