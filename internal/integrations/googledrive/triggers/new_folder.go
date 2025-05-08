package triggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type newFolderTriggerProps struct {
	newFileTriggerProps
}

type NewFolderTrigger struct {
	timezoneOptions []*sdkcore.AutoFormSchema
	hodOptions      []*sdkcore.AutoFormSchema
}

func (t *NewFolderTrigger) Name() string {
	return "New Folder"
}

func (t *NewFolderTrigger) Description() string {
	return "New Folder trigger is designed to monitor a specific folder or directory for new files or subfolders. Whenever a new folder is created within the monitored directory, this trigger will automatically initiate the workflow automation process, allowing you to streamline tasks and automate workflows related to file organization, data processing, or other business-critical activities."
}

func (t *NewFolderTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewFolderTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newFolderDocs,
	}
}

func (t *NewFolderTrigger) Icon() *string {
	return nil
}

func (t *NewFolderTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"parentFolder":      shared.RegisterParentFoldersProp(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

// Start initializes the newFolderTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewFolderTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newFolderTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewFolderTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newFolderTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewFolderTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newFolderTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	fmt.Printf("ctx.Metadata.LastRun %+v \n", ctx.Metadata().LastRun)

	var qarr []string
	if input.ParentFolder != nil {
		qarr = append(qarr, fmt.Sprintf("'%v' in parents", *input.ParentFolder))
	}

	if input.CreatedTime == nil {
		input.CreatedTime = ctx.Metadata().LastRun
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

func (t *NewFolderTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewFolderTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewFolderTrigger() sdk.Trigger {
	return &NewFolderTrigger{}
}
