package actions

import (
	"bytes"
	"context"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type createFileActionProps struct {
	FileName          string  `json:"fileName"`
	Content           string  `json:"content"`
	ContentType       string  `json:"contentType"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type CreateFileAction struct{}

// Metadata returns metadata about the action
func (a *CreateFileAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_file",
		DisplayName:   "Create File",
		Description:   "Creates a new file with a specified name and content, allowing you to store and manage data within your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: createFileDocs,
		SampleOutput: map[string]any{
			"kind":     "drive#file",
			"mimeType": "text/plain",
			"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
			"name":     "example.txt",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateFileAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_file", "Create File")

	form.TextField("fileName", "File name").
		Placeholder("example.txt").
		Required(true).
		HelpText("The name of the new file with extension.")

	form.TextareaField("content", "Content").
		Placeholder("Enter file content here...").
		Required(true).
		HelpText("The content to be written to the file.")

	form.SelectField("contentType", "Content Type").
		AddOption("text", "Text").
		AddOption("csv", "CSV").
		AddOption("xml", "XML").
		DefaultValue("text").
		Required(true).
		HelpText("The type of content being created.")

	shared.RegisterParentFoldersProp(form)

	form.CheckboxField("includeTeamDrives", "Include Team Drives").
		Required(false).
		HelpText("Whether to include team drives in the folder selection.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateFileAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateFileAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createFileActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// Determine MIME type based on content type
	var mimeType string
	switch input.ContentType {
	case "csv":
		mimeType = "text/csv"
	case "xml":
		mimeType = "text/xml"
	default:
		mimeType = "text/plain"
	}

	// Create file metadata
	fileMetadata := &drive.File{
		Name:     input.FileName,
		MimeType: mimeType,
	}

	// Add parent folder if specified
	if input.ParentFolder != nil {
		fileMetadata.Parents = []string{*input.ParentFolder}
	}

	// Create the file with content
	file, err := driveService.Files.Create(fileMetadata).
		Media(bytes.NewReader([]byte(input.Content)), googleapi.ContentType(mimeType)).
		Fields("id, name, mimeType, webViewLink, kind, createdTime").
		SupportsAllDrives(input.IncludeTeamDrives).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return file, nil
}

func NewCreateFileAction() sdk.Action {
	return &CreateFileAction{}
}
