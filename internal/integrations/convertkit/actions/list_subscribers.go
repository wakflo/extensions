package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listSubscribersActionProps struct {
	Page     int  `json:"page"`
	PageSize int  `json:"page-size"`
	FromDate *int `json:"from-date"`
}

type ListSubscribersAction struct{}

func (a *ListSubscribersAction) Name() string {
	return "List Subscribers"
}

func (a *ListSubscribersAction) Description() string {
	return "Retrieve a list of subscribers from your ConvertKit account with their details and tags."
}

func (a *ListSubscribersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListSubscribersAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listSubscribersDocs,
	}
}

func (a *ListSubscribersAction) Icon() *string {
	icon := "mdi:account-group"
	return &icon
}

func (a *ListSubscribersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("The page number to retrieve (pagination)").
			Build(),
	}
}

func (a *ListSubscribersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[listSubscribersActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	path := "/subscribers?api_secret=" + ctx.Auth.Extra["api-secret"]

	response, err := shared.GetConvertKitClient(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, err
	}

	subscribers, ok := responseMap["subscribers"]
	if !ok {
		return nil, fmt.Errorf("failed to extract subscribers from response")
	}

	return subscribers, nil
}

func (a *ListSubscribersAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListSubscribersAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"subscribers": []map[string]any{
			{
				"id":            "123456",
				"first_name":    "Jane",
				"email_address": "jane@example.com",
				"state":         "active",
				"created_at":    "2023-01-15T10:30:00Z",
				"fields": map[string]string{
					"company": "Acme Inc",
				},
			},
			{
				"id":            "789012",
				"first_name":    "John",
				"email_address": "john@example.com",
				"state":         "active",
				"created_at":    "2023-01-16T14:20:00Z",
				"fields": map[string]string{
					"company": "XYZ Corp",
				},
			},
		},
		"total_subscribers": "156",
		"page":              "1",
		"page_size":         "50",
	}
}

func (a *ListSubscribersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListSubscribersAction() sdk.Action {
	return &ListSubscribersAction{}
}
