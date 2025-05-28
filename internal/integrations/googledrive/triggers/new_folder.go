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

type newFolderTriggerProps struct {
	newFileTriggerProps
}

type NewFolderTrigger struct{}

func (t *NewFolderTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_folder",
		DisplayName:   "New Folder",
		Description:   "Triggers when a new folder is created or uploaded to a specified directory or location, allowing you to automate workflows and processes as soon as a new folder becomes available.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newFolderDocs,
		SampleOutput: map[string]any{
			"files": []map[string]any{
				{
					"kind": "drive#file",
				},
			},
		},
	}
}

func (t *NewFolderTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("google-drive-new-folder", "New File")

	shared.RegisterParentFoldersProp(form)

	// Add include team drives field
	form.CheckboxField("includeTeamDrives", "Include Team Drives").
		Placeholder("Enter a value for Include Team Drives.").
		Required(false).
		HelpText("Whether to include team drives in the folder selection.")

	schema := form.Build()

	return schema
}

func (t *NewFolderTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewFolderTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

// Start initializes the newFolderTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewFolderTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newFolderTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewFolderTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newFolderTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewFolderTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newFolderTriggerProps](ctx)
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
		lr := ctx.LastRun()
		input.CreatedTime = lr
	}

	if input.CreatedTime != nil {
		op := ">"
		if input.CreatedTimeOp != nil {
			op = *input.CreatedTimeOp
		}

		qarr = append(qarr, fmt.Sprintf(`createdTime %v '%v'`, op, input.CreatedTime.UTC().Format(time.RFC3339)))
	}

	qarr = append(qarr, "trashed = false")
	q := fmt.Sprintf("%v %v", "mimeType='application/vnd.google-apps.folder'  and ", strings.Join(qarr, " and "))

	req := driveService.Files.List().
		IncludeItemsFromAllDrives(input.IncludeTeamDrives).
		SupportsAllDrives(input.IncludeTeamDrives).
		Fields("files(id, name, mimeType, webViewLink, kind, createdTime)").
		Q(q)

	files, err := req.Do()
	if err != nil {
		return nil, err
	}

	return files.Files, nil
}

func (t *NewFolderTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewFolderTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewFolderTrigger() sdk.Trigger {
	return &NewFolderTrigger{}
}
