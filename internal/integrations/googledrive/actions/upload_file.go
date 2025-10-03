package actions

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gabriel-vasile/mimetype"
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

	form.TextField("file", "File URL").
		Placeholder("Enter the file URL or base64 data").
		Required(true).
		HelpText("The file URL or base64 encoded data to upload")

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

// processFileString handles both URL and base64 inputs
func (a *UploadFileAction) processFileString(fileStr string) (io.Reader, string, error) {
	// Check if it's a URL
	if strings.HasPrefix(fileStr, "http://") || strings.HasPrefix(fileStr, "https://") {
		// Download file from URL
		resp, err := http.Get(fileStr)
		if err != nil {
			return nil, "", fmt.Errorf("failed to download file from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, "", fmt.Errorf("failed to download file: HTTP status %d", resp.StatusCode)
		}

		// Read the entire response body
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read response body: %w", err)
		}

		// Detect MIME type
		mtype := mimetype.Detect(data)

		return bytes.NewReader(data), mtype.String(), nil
	}

	// Check if it's a data URI (e.g., data:image/png;base64,...)
	if strings.HasPrefix(fileStr, "data:") {
		parts := strings.SplitN(fileStr, ",", 2)
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("invalid data URI format")
		}

		// Extract MIME type from the prefix
		prefix := parts[0]
		var mimeType string
		if strings.Contains(prefix, ";") {
			mimeType = strings.TrimPrefix(strings.Split(prefix, ";")[0], "data:")
		}

		// Decode base64 data
		data, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode base64 data: %w", err)
		}

		// If no MIME type was specified, detect it
		if mimeType == "" {
			mtype := mimetype.Detect(data)
			mimeType = mtype.String()
		}

		return bytes.NewReader(data), mimeType, nil
	}

	// Try to decode as plain base64
	data, err := base64.StdEncoding.DecodeString(fileStr)
	if err != nil {
		// Try URL-safe base64 encoding
		data, err = base64.URLEncoding.DecodeString(fileStr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid file string: not a valid URL, data URI, or base64 string")
		}
	}

	// Detect MIME type
	mtype := mimetype.Detect(data)

	return bytes.NewReader(data), mtype.String(), nil
}

// Perform executes the action with the given context and input
func (a *UploadFileAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[uploadFileActionProps](ctx)
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

	// Process the file string
	fileReader, mimeType, err := a.processFileString(input.File)
	if err != nil {
		return nil, err
	}

	var parents []string
	if input.ParentFolder != nil {
		parents = append(parents, *input.ParentFolder)
	}

	in := &drive.File{
		MimeType: mimeType,
		Name:     input.FileName,
		Parents:  parents,
	}

	result, err := driveService.Files.Create(in).
		Media(fileReader).
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
