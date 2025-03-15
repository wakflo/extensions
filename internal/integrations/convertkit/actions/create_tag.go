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

type createTagActionProps struct {
	TagName string `json:"tag_name"`
}

type CreateTagAction struct{}

func (a *CreateTagAction) Name() string {
	return "Create Tag"
}

func (a *CreateTagAction) Description() string {
	return "Create a new tag in your ConvertKit account."
}

func (a *CreateTagAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateTagAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createTagDocs,
	}
}

func (a *CreateTagAction) Icon() *string {
	icon := "mdi:tag"
	return &icon
}

func (a *CreateTagAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tag_name": autoform.NewShortTextField().
			SetDisplayName("Tag Name").
			SetDescription("Name of the tag to create").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateTagAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTagActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"api_secret": ctx.Auth.Extra["api-secret"],
		"tag": map[string]string{
			"name": input.TagName,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	response, err := shared.GetConvertKitClient("/tags", http.MethodPost, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *CreateTagAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateTagAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"tag": map[string]any{
			"id":         "123456",
			"name":       "Example Tag",
			"created_at": "2024-03-15T10:30:00Z",
		},
	}
}

func (a *CreateTagAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateTagAction() sdk.Action {
	return &CreateTagAction{}
}
