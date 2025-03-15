package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/convertkit/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type ListTagsAction struct{}

func (a *ListTagsAction) Name() string {
	return "List Tags"
}

func (a *ListTagsAction) Description() string {
	return "Retrieve a list of all tags in your ConvertKit account."
}

func (a *ListTagsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListTagsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listTagsDocs,
	}
}

func (a *ListTagsAction) Icon() *string {
	icon := "mdi:tag-multiple"
	return &icon
}

func (a *ListTagsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of subscribers to retrieve (default: 50)").
			SetDefaultValue(50).
			Build(),
	}
}

func (a *ListTagsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	// Construct the path with API key
	path := "/tags?api_key=" + ctx.Auth.Extra["api-key"]

	// Make the API request
	response, err := shared.GetConvertKitClient(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	// Return the response
	return response, nil
}

func (a *ListTagsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListTagsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"tags": []map[string]any{
			{
				"id":         "1001",
				"name":       "Newsletter",
				"created_at": "2023-05-15T10:30:00Z",
			},
			{
				"id":         "1002",
				"name":       "Product Updates",
				"created_at": "2023-06-20T14:45:00Z",
			},
			{
				"id":         "1003",
				"name":       "New Customers",
				"created_at": "2023-07-12T09:15:00Z",
			},
		},
	}
}

func (a *ListTagsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListTagsAction() sdk.Action {
	return &ListTagsAction{}
}
