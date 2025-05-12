package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listSubscribersActionProps struct {
	ListID string `json:"listId"`
	Page   int    `json:"page"`
}

type ListSubscribersAction struct{}

// Metadata returns metadata about the action
func (a *ListSubscribersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_subscribers",
		DisplayName:   "List Active Subscribers",
		Description:   "Retrieve a list of active subscribers from a specific list.",
		Type:          core.ActionTypeAction,
		Documentation: listSubscribersDocs,
		Icon:          "mdi:account-group",
		SampleOutput: map[string]interface{}{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListSubscribersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_subscribers", "List Active Subscribers")

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("listId", "List").
	//	Placeholder("Select a list").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The list to retrieve subscribers from.")

	form.NumberField("page", "Page").
		Placeholder("Enter page number").
		Required(false).
		HelpText("The page number to retrieve (for pagination).")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListSubscribersAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListSubscribersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listSubscribersActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
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
		authCtx.Extra["api-key"],
		authCtx.Extra["client-id"],
		endpoint,
		http.MethodGet,
		nil)
	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

func NewListSubscribersAction() sdk.Action {
	return &ListSubscribersAction{}
}
