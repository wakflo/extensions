package triggers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type newFileProps struct {
	ParentFolder       *string    `json:"parentFolder"`
	IncludeTeamDrives  bool       `json:"includeTeamDrives"`
	IncludeFileContent bool       `json:"includeFileContent"`
	CreatedTime        *time.Time `json:"createdTime"`
	CreatedTimeOp      *string    `json:"createdTimeOp"`
}

type NewFileTrigger struct {
	getParentFolders func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error)
}

func (e *NewFileTrigger) Name() string {
	return "New File"
}

func (e *NewFileTrigger) Description() string {
	return "Triggers a workflow to by providing seamless automation for monitoring newly created files with customizable options like folder selection and additional file details."
}

func (e *NewFileTrigger) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &newFileDocs,
	}
}

func (e *NewFileTrigger) Icon() *string {
	return nil
}

func (e *NewFileTrigger) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"kind":     "drive#file",
		"mimeType": "image/jpeg",
		"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
		"name":     "example.jpg",
	}
}

func (e *NewFileTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"parentFolder": autoform.NewDynamicField(sdkcore.String).
			SetDisplayName("Parent Folder").
			SetDescription("select parent folder").
			SetDynamicOptions(&e.getParentFolders).
			SetDependsOn([]string{"connection"}).
			SetRequired(false).Build(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
		"includeFileContent": autoform.NewBooleanField().
			SetDisplayName("Include File Content").
			SetDescription("Include the file content in the output. This will increase the time taken to fetch the files and might cause issues with large files.").
			SetDefaultValue(false).
			Build(),
	}
}

func (e *NewFileTrigger) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (e *NewFileTrigger) Start(ctx integration.LifecycleContext) error {
	return nil
}

func (e *NewFileTrigger) Stop(ctx integration.LifecycleContext) error {
	return nil
}

func (e *NewFileTrigger) Execute(ctx integration.ExecuteContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[newFileProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

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
		return shared.HandleFileContent(&ctx.BaseContext, files.Files, driveService)
	}
	return files.Files, nil
}

func (e *NewFileTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypeScheduled
}

func (e *NewFileTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
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

func NewNewFileTrigger() integration.Trigger {
	getParentFolders := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := fastshot.NewClient("https://www.googleapis.com/drive/v3").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/files").Query().
			AddParams(map[string]string{
				"q": "mimeType='application/vnd.google-apps.folder' and trashed = false",
				/*"supportsTeamDrives": "true",
				"supportsAllDrives":  "true",*/
			}).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body shared.ListFileResponse
		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(body.Files, len(body.Files))
	}

	return &NewFileTrigger{
		getParentFolders: getParentFolders,
	}
}
