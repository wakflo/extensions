package triggers

import (
	"context"
	"errors"
	"log"

	"github.com/wakflo/extensions/internal/integrations/toggl/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newProjectProps struct {
	WorkspaceID string `json:"workspace_id"`
}

type NewProjectTrigger struct {
	timezoneOptions []*sdkcore.AutoFormSchema
	hodOptions      []*sdkcore.AutoFormSchema
}

func (e *NewProjectTrigger) Name() string {
	return "New Project"
}

func (e *NewProjectTrigger) Description() string {
	return "Schedules a workflow to run every hour"
}

func (e *NewProjectTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newProjectDocs,
	}
}

func (e *NewProjectTrigger) Icon() *string {
	return nil
}

func (e *NewProjectTrigger) SampleData() sdkcore.JSON {
	return nil
}

func (e *NewProjectTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace_id": shared.GetWorkSpaceInput(),
	}
}

func (e *NewProjectTrigger) Auth() *sdk.Auth {
	return nil
}

func (e *NewProjectTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (e *NewProjectTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (e *NewProjectTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing toggl api key")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	input, err := sdk.InputToTypeSafely[newProjectProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun

	var updatedTime int64
	if lastRunTime != nil {
		updatedTime = lastRunTime.UTC().Unix()
	}

	response, err := shared.GetProject(apiKey, input.WorkspaceID, updatedTime)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func (e *NewProjectTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypeScheduled
}

func (e *NewProjectTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
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

func NewNewProjectTrigger() sdk.Trigger {
	return &NewProjectTrigger{}
}
