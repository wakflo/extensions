package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listSubscribersActionProps struct {
	ListID string `json:"listId"`
	Page   int    `json:"page"`
}

type ListSubscribersAction struct{}

func (a *ListSubscribersAction) Name() string {
	return "List Active Subscribers"
}

func (a *ListSubscribersAction) Description() string {
	return "Retrieve a list of active subscribers from a specific list."
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
		"listId": shared.GetCreateSendSubscriberListsInput(),
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("The page number to retrieve (for pagination).").
			SetRequired(false).
			Build(),
	}
}

func (a *ListSubscribersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listSubscribersActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("lists/%s/active.json", input.ListID)

	queryParams := make(map[string]string)

	page := input.Page
	if page <= 0 {
		page = 1
	}
	queryParams["page"] = strconv.Itoa(page)

	if len(queryParams) > 0 {
		endpoint += shared.BuildQueryString(queryParams)
	}

	subscribers, err := shared.GetCampaignMonitorClient(
		ctx.Auth.Extra["api-key"],
		ctx.Auth.Extra["client-id"],
		endpoint,
		http.MethodGet,
		nil)
	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

func (a *ListSubscribersAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListSubscribersAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"Results": []interface{}{
			map[string]interface{}{
				"EmailAddress": "subscriber1@example.com",
				"Name":         "John Doe",
				"Date":         "2023-07-10T14:30:00",
				"State":        "Active",
				"CustomFields": []interface{}{
					map[string]interface{}{
						"Key":   "City",
						"Value": "New York",
					},
					map[string]interface{}{
						"Key":   "Age",
						"Value": "34",
					},
				},
			},
			map[string]interface{}{
				"EmailAddress": "subscriber2@example.com",
				"Name":         "Jane Smith",
				"Date":         "2023-06-20T09:15:00",
				"State":        "Active",
				"CustomFields": []interface{}{
					map[string]interface{}{
						"Key":   "City",
						"Value": "Los Angeles",
					},
					map[string]interface{}{
						"Key":   "Age",
						"Value": "29",
					},
				},
			},
		},
		"ResultsOrderedBy":     "email",
		"OrderDirection":       "asc",
		"PageNumber":           "1",
		"PageSize":             "1000",
		"RecordsOnThisPage":    "2",
		"TotalNumberOfRecords": "2",
		"NumberOfPages":        "1",
	}
}

func (a *ListSubscribersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListSubscribersAction() sdk.Action {
	return &ListSubscribersAction{}
}
