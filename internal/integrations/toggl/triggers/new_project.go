package triggers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/toggl/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type newProjectProps struct {
	WorkspaceID string `json:"workspace_id"`
}

type NewProjectTrigger struct{}

func (e *NewProjectTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_project",
		DisplayName:   "New Project",
		Description:   "Schedules a workflow to run every hour",
		Type:          core.TriggerTypeScheduled,
		Documentation: newProjectDocs,
		SampleOutput:  nil,
	}
}

func (e *NewProjectTrigger) Auth() *core.AuthMetadata {
	return nil
}

func (e *NewProjectTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new-project", "New Project")

	shared.RegisterWorkspacesProp(form)

	schema := form.Build()
	return schema
}

func (e *NewProjectTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (e *NewProjectTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (e *NewProjectTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing toggl api key")
	}
	apiKey := authCtx.Extra["api-key"]

	input, err := sdk.InputToTypeSafely[newProjectProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the last run time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	var updatedTime int64
	if lastRun != nil {
		lastRunTime, ok := lastRun.(*time.Time)
		if ok && lastRunTime != nil {
			updatedTime = lastRunTime.UTC().Unix()
		}
	}

	response, err := shared.GetProject(apiKey, input.WorkspaceID, updatedTime)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func (e *NewProjectTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{
		Schedule: &core.ScheduleTriggerCriteria{
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
