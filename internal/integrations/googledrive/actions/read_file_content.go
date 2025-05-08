package actions

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"slices"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type readFileContentActionProps struct {
	FileID   string  `json:"fileId"`
	FileName *string `json:"fileName"`
}

type ReadFileContentAction struct{}

// Metadata returns metadata about the action
func (a *ReadFileContentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "read_file_content",
		DisplayName:   "Read File Content",
		Description:   "Reads the content of a specified file and returns it as a string or binary data, depending on the file type. This action is useful when you need to extract information from a file or process its contents in your workflow automation.",
		Type:          core.ActionTypeAction,
		Documentation: readFileContentDocs,
		SampleOutput: map[string]any{
			"fileData": "https://example.com/file.txt",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ReadFileContentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("read_file_content", "Read File Content")

	form.TextField("fileId", "File ID").
		Placeholder("Enter a file ID").
		Required(true).
		HelpText("File ID coming from | New File -> id |")

	form.TextField("fileName", "File Name").
		Placeholder("Enter a file name").
		Required(true).
		HelpText("Destination File name")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ReadFileContentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ReadFileContentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[readFileContentActionProps](ctx)
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

	// Since we can't directly use the shared.DownloadFile function with the new SDK context,
	// we'll implement the file download logic here
	file, err := driveService.Files.Get(input.FileID).Do()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	var rsp *http.Response
	if slices.Contains([]string{
		"application/vnd.google-apps.document",
		"application/vnd.google-apps.spreadsheet",
		"application/vnd.google-apps.presentation",
	}, file.MimeType) {
		rsp, err = driveService.Files.Export(file.Id, file.MimeType).Download()
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()
		_, err = buf.ReadFrom(rsp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		rsp, err = driveService.Files.Get(file.Id).Download()
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()
		_, err = buf.ReadFrom(rsp.Body)
		if err != nil {
			return nil, err
		}
	}

	ext, err := mime.ExtensionsByType(file.MimeType)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s.%s", file.Name, ext)
	if input.FileName != nil {
		name = fmt.Sprintf("%s.%s", *input.FileName, ext)
	}

	// Encode the file content as base64
	fileData := base64.StdEncoding.EncodeToString(buf.Bytes())

	return map[string]interface{}{
		"fileName": name,
		"fileData": fileData,
		"mimeType": file.MimeType,
	}, nil
}

func NewReadFileContentAction() sdk.Action {
	return &ReadFileContentAction{}
}
