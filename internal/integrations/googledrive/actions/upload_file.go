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

type uploadFileActionProps struct {
	FileName          string  `json:"fileName"`
	File              string  `json:"file"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type UploadFileAction struct{}

// Metadata returns metadata about the action
func (a *UploadFileAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "upload_file",
		DisplayName:   "Upload File",
		Description:   "Upload File: This integration action allows you to upload files from various sources such as cloud storage services, local file systems, or email attachments to your workflow. You can specify the file type, size limit, and other parameters to control the upload process. The uploaded file is then stored in a designated location within your workflow, making it easily accessible for further processing or analysis.",
		Type:          core.ActionTypeAction,
		Documentation: uploadFileDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UploadFileAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("upload_file", "Upload File")

	form.TextField("fileName", "Name").
		Placeholder("Enter a file name").
		Required(true).
		HelpText("The name of the new file")

	form.FileField("file", "File").
		Placeholder("Select a file to upload").
		Required(true).
		HelpText("The file URL or base64 to upload")

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
func (a *UploadFileAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UploadFileAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[uploadFileActionProps](ctx)
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

	result, err := driveService.Files.Create(in).
		Media(fileData.Data).
		SupportsAllDrives(input.IncludeTeamDrives).
		Do()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewUploadFileAction() sdk.Action {
	return &UploadFileAction{}
}
