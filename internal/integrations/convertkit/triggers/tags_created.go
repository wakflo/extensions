package triggers

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type tagCreatedTriggerProps struct {
	Limit int `json:"limit"`
}

type TagCreatedTrigger struct{}

func (t *TagCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "tag_created",
		DisplayName:   "Tag Created",
		Description:   "Triggers a workflow when a new tag is created in your ConvertKit account, allowing you to automate follow-up actions.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: TagCreatedDocs,
		Icon:          "mdi:tag-plus",
		SampleOutput: map[string]any{
			"tags": []map[string]any{
				{
					"id":         "12345",
					"name":       "Newsletter Subscribers",
					"created_at": "2024-03-15T10:30:00Z",
				},
				{
					"id":         "12346",
					"name":       "Product Interest",
					"created_at": "2024-03-15T10:35:00Z",
				},
			},
			"total_tags": "2",
		},
	}
}

func (t *TagCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *TagCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TagCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("convertkit-tag-created", "Tag Created")

	form.NumberField("limit", "limit").
		Placeholder("Limit").
		HelpText("Maximum number of tags to retrieve (default: 50)").
		DefaultValue(50).
		Required(false)

	schema := form.Build()

	return schema
}

func (t *TagCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TagCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *TagCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[tagCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	path := "/tags?api_key=" + authCtx.Extra["api-key"]

	response, err := shared.GetConvertKitClient(path, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags: %v", err)
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, err
	}

	tagsData, ok := responseMap["tags"]
	if !ok {
		return nil, errors.New("failed to extract tags from response")
	}

	return tagsData, nil
}

func (t *TagCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func NewTagCreatedTrigger() sdk.Trigger {
	return &TagCreatedTrigger{}
}
