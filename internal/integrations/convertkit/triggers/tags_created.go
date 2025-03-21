package triggers

import (
	"context"
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type tagCreatedTriggerProps struct {
	Limit int `json:"limit"`
}

type TagCreatedTrigger struct{}

func (t *TagCreatedTrigger) Name() string {
	return "Tag Created"
}

func (t *TagCreatedTrigger) Description() string {
	return "Triggers a workflow when a new tag is created in your ConvertKit account, allowing you to automate follow-up actions."
}

func (t *TagCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *TagCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &TagCreatedDocs,
	}
}

func (t *TagCreatedTrigger) Icon() *string {
	icon := "mdi:tag-plus"
	return &icon
}

func (t *TagCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of tags to retrieve (default: 50)").
			Build(),
	}
}

func (t *TagCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TagCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *TagCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[tagCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	path := "/tags?api_key=" + ctx.Auth.Extra["api-key"]

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

func (t *TagCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *TagCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func NewTagCreatedTrigger() sdk.Trigger {
	return &TagCreatedTrigger{}
}
