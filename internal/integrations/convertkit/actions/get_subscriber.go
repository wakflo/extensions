package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getSubscriberActionProps struct {
	SubscriberID int `json:"subscriber_id"`
}

type GetSubscriberAction struct{}

func (a *GetSubscriberAction) Name() string {
	return "Get Subscriber"
}

func (a *GetSubscriberAction) Description() string {
	return "Retrieve detailed information about a specific subscriber by ID."
}

func (a *GetSubscriberAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetSubscriberAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getSubscriberDocs,
	}
}

func (a *GetSubscriberAction) Icon() *string {
	icon := "mdi:account-details"
	return &icon
}

func (a *GetSubscriberAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"subscriber_id": autoform.NewNumberField().
			SetDisplayName("Subscriber ID").
			SetDescription("ID of the subscriber to retrieve").
			SetRequired(true).
			Build(),
	}
}

func (a *GetSubscriberAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getSubscriberActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/subscribers/%d?api_secret=%s", input.SubscriberID, ctx.Auth.Extra["api-secret"])

	response, err := shared.GetConvertKitClient(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *GetSubscriberAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetSubscriberAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"subscriber": map[string]any{
			"id":            "12345",
			"first_name":    "Jon",
			"email_address": "jon.snow@example.com",
			"state":         "active",
			"created_at":    "2023-05-15T10:30:00Z",
			"fields": map[string]any{
				"last_name": "Snow",
				"company":   "Night's Watch",
				"location":  "The Wall",
			},
		},
	}
}

func (a *GetSubscriberAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetSubscriberAction() sdk.Action {
	return &GetSubscriberAction{}
}
