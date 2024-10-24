package notion

import (
	"errors"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type newPageTriggerProps struct {
	DatabaseID    string `json:"database"`
	CheckInterval int    `json:"check_interval"`
}

type TriggerNewPageCreated struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewPageCreated() *TriggerNewPageCreated {
	return &TriggerNewPageCreated{
		options: &sdk.TriggerInfo{
			Name:        "New Page Created",
			Description: "Triggers workflow when a new page is created in a Notion database",
			RequireAuth: true,
			Auth:        sharedAuth, // Assume this is defined elsewhere
			Input: map[string]*sdkcore.AutoFormSchema{
				"database": getNotionDatabasesInput("Database", "Select a database", true),
				"check_interval": autoform.NewNumberField().
					SetDisplayName("Check Interval (minutes)").
					SetDescription("How often to check for new pages (in minutes).").
					SetRequired(true).
					SetDefaultValue(5).
					Build(),
			},
			Type:     sdkcore.TriggerTypeCron,
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewPageCreated) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[newPageTriggerProps](ctx)

	if input.DatabaseID == "" {
		return nil, errors.New("database_id is required")
	}

	if input.CheckInterval <= 0 {
		return nil, errors.New("check_interval must be greater than 0")
	}

	now := time.Now().UTC()
	var lastChecked time.Time

	if ctx.Metadata.LastRun != nil {
		lastChecked = *ctx.Metadata.LastRun
	} else {
		lastChecked = now.Add(-time.Duration(input.CheckInterval) * time.Minute)
	}

	newPages, err := queryNewPages(accessToken, input.DatabaseID, lastChecked)
	if err != nil {
		return nil, err
	}

	if len(newPages) > 0 {
		return map[string]interface{}{"new_pages": newPages}, nil
	}

	return map[string]interface{}{"message": "No new pages found"}, nil
}

func (t *TriggerNewPageCreated) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewPageCreated) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewPageCreated) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewPageCreated) GetInfo() *sdk.TriggerInfo {
	return t.options
}
