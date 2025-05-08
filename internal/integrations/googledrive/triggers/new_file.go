package triggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type newFileTriggerProps struct {
	ParentFolder       *string    `json:"parentFolder"`
	IncludeTeamDrives  bool       `json:"includeTeamDrives"`
	IncludeFileContent bool       `json:"includeFileContent"`
	CreatedTime        *time.Time `json:"createdTime"`
	CreatedTimeOp      *string    `json:"createdTimeOp"`
}

type NewFileTrigger struct {
}

func (t *NewFileTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_file",
		DisplayName:   "New File",
		Description:   "Triggers when a new file is created or uploaded to a specified directory or location, allowing you to automate workflows and processes as soon as a new file becomes available.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newFileDocs,
		SampleOutput: map[string]any{
			"files": []map[string]any{
				{
					"kind": "drive#file",
				},
			},
		},
	}
}

func (t *NewFileTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewFileTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewFileTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("google-drive-new-file", "New File")

	shared.RegisterParentFoldersProp(form)

	form.CheckboxField("includeFileContent", "Include File Content").
		Placeholder("Enter a value for Include File Content.").
		Required(false).
		HelpText("Include the file content in the output. This will increase the time taken to fetch the files and might cause issues with large files.")

	// Add include team drives field
	form.CheckboxField("includeTeamDrives", "Include Team Drives").
		Placeholder("Enter a value for Include Team Drives.").
		Required(false).
		HelpText("Whether to include team drives in the folder selection.")

	schema := form.Build()

	return schema
}

// Start initializes the newFileTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewFileTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newFileTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewFileTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newFileTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewFileTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newFileTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth().TokenSource))
	if err != nil {
		return nil, err
	}

	var qarr []string
	if input.ParentFolder != nil {
		qarr = append(qarr, fmt.Sprintf("'%v' in parents", *input.ParentFolder))
	}
	if input.CreatedTime == nil {
		lr, err := ctx.GetMetadata("lastRun")
		if err != nil {
			return nil, err
		}
		input.CreatedTime = lr.(*time.Time)
	}
	if input.CreatedTime != nil {
		op := ">"
		if input.CreatedTimeOp != nil {
			op = *input.CreatedTimeOp
		}
		qarr = append(qarr, fmt.Sprintf(`createdTime %v '%v'`, op, input.CreatedTime.UTC().Format(time.RFC3339)))
	}

	qarr = append(qarr, "trashed = false")
	q := fmt.Sprintf("%v %v", "mimeType!='application/vnd.google-apps.folder'  and ", strings.Join(qarr, " and "))

	req := driveService.Files.List().
		IncludeItemsFromAllDrives(input.IncludeTeamDrives).
		SupportsAllDrives(input.IncludeTeamDrives).
		Fields("files(id, name, mimeType, webViewLink, kind, createdTime)").
		Q(q)

	files, err := req.Do()
	if err != nil {
		return nil, err
	}

	if input.IncludeFileContent {
		return shared.HandleFileContent(ctx, files.Files, driveService)
	}
	return files.Files, nil
}

func (t *NewFileTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewFileTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"kind":     "drive#file",
		"mimeType": "image/jpeg",
		"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
		"name":     "example.jpg",
	}
}

func NewNewFileTrigger() sdk.Trigger {
	return &NewFileTrigger{}
}
