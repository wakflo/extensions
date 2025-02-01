package triggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type newFolderProps = newFileProps

type NewFolderTrigger struct {
	timezoneOptions []*sdkcore.AutoFormSchema
	hodOptions      []*sdkcore.AutoFormSchema
}

func (e *NewFolderTrigger) Name() string {
	return "New Folder"
}

func (e *NewFolderTrigger) Description() string {
	return "Schedules a workflow to run every hour"
}

func (e *NewFolderTrigger) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &newFolderDocs,
	}
}

func (e *NewFolderTrigger) Icon() *string {
	return nil
}

func (e *NewFolderTrigger) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (e *NewFolderTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"parentFolder":      shared.GetParentFoldersInput(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (e *NewFolderTrigger) Auth() *integration.Auth {
	return nil
}

func (e *NewFolderTrigger) Start(ctx integration.LifecycleContext) error {
	return nil
}

func (e *NewFolderTrigger) Stop(ctx integration.LifecycleContext) error {
	return nil
}

func (e *NewFolderTrigger) Execute(ctx integration.ExecuteContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[newFolderProps](ctx.BaseContext)
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

func (e *NewFolderTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypeScheduled
}

func (e *NewFolderTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{
		Schedule: &sdkcore.ScheduleTriggerCriteria{
			CronExpression: "",
			StartTime:      nil,
			EndTime:        nil,
			TimeZone:       "",
			Enabled:        true,
		},
	}
}

func NewNewFolderTrigger() integration.Trigger {
	return &NewFolderTrigger{}
}
