package actions

import (
	"context"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type getFileActionProps struct {
	FileID string `json:"fileId"`
}

type GetFileAction struct{}

// Metadata returns metadata about the action
func (a *GetFileAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_file",
		DisplayName:   "Get File",
		Description:   "Retrieves a file from a specified location and makes it available for further processing in the workflow.",
		Type:          core.ActionTypeAction,
		Documentation: getFileDocs,
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
func (a *GetFileAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_file", "Get File")

	form.TextField("fileId", "File / Folder ID").
		Placeholder("Enter a file or folder ID").
		Required(true).
		HelpText("The ID of the file/folder to search for.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetFileAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetFileAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getFileActionProps](ctx)
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

	file, err := driveService.Files.Get(input.FileID).
		Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()

	if err != nil {
		return nil, err
	}

	return file, nil
}

func NewGetFileAction() sdk.Action {
	return &GetFileAction{}
}
