package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type tagSubscriberActionProps struct {
	TagID int    `json:"tag_id"`
	Email string `json:"email"`
}

type TagSubscriberAction struct{}

func (a *TagSubscriberAction) Name() string {
	return "Tag Subscriber"
}

func (a *TagSubscriberAction) Description() string {
	return "Apply a tag to a subscriber in your ConvertKit account."
}

func (a *TagSubscriberAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *TagSubscriberAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &tagSubscriberDocs,
	}
}

func (a *TagSubscriberAction) Icon() *string {
	icon := "mdi:tag-plus"
	return &icon
}

func (a *TagSubscriberAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tag_id": autoform.NewNumberField().
			SetDisplayName("Tag ID").
			SetDescription("ID of the tag to apply").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Email address of the subscriber").
			SetRequired(true).
			Build(),
	}
}

func (a *TagSubscriberAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[tagSubscriberActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"api_secret": ctx.Auth.Extra["api-secret"],
		"email":      input.Email,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	path := fmt.Sprintf("/tags/%d/subscribe", input.TagID)

	response, err := shared.GetConvertKitClient(path, http.MethodPost, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TagSubscriberAction) Auth() *sdk.Auth {
	return nil
}

func (a *TagSubscriberAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"subscription": map[string]any{
			"id":                "12345",
			"state":             "inactive",
			"created_at":        "2024-03-15T10:30:00Z",
			"source":            nil,
			"referrer":          nil,
			"subscribable_id":   "789",
			"subscribable_type": "tag",
			"subscriber": map[string]any{
				"id": "54321",
			},
		},
	}
}

func (a *TagSubscriberAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewTagSubscriberAction() sdk.Action {
	return &TagSubscriberAction{}
}
