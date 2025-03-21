package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type GetSubscriberListsAction struct{}

func (a *GetSubscriberListsAction) Name() string {
	return "Get Subscriber Lists"
}

func (a *GetSubscriberListsAction) Description() string {
	return "Retrieve all subscriber lists for a client."
}

func (a *GetSubscriberListsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetSubscriberListsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getSubscriberList,
	}
}

func (a *GetSubscriberListsAction) Icon() *string {
	icon := "mdi:format-list-bulleted"
	return &icon
}

func (a *GetSubscriberListsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("The page number to retrieve (for pagination).").
			SetRequired(false).
			Build(),
	}
}

func (a *GetSubscriberListsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	clientID, clientIDExists := ctx.Auth.Extra["client-id"]
	apiKey, apiKeyExists := ctx.Auth.Extra["api-key"]

	if !clientIDExists || !apiKeyExists || clientID == "" || apiKey == "" {
		return nil, errors.New("client ID and API Key are required")
	}

	// Construct API URL
	apiURL := fmt.Sprintf("https://api.createsend.com/api/v3.3/clients/%s/lists.json", clientID)

	// Create request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(apiKey, "")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s", body)
	}

	var lists []map[string]interface{}
	if err := json.Unmarshal(body, &lists); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return lists, nil
}

func (a *GetSubscriberListsAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetSubscriberListsAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"lists": []map[string]interface{}{
			{
				"id":   "a58ee1d3039b8bec838e6d1482a8a965",
				"name": "List One",
			},
			{
				"id":   "99bc35084a5739127a8ab81eae5bd305",
				"name": "List Two",
			},
		},
		"count": "2",
	}
}

func (a *GetSubscriberListsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetSubscriberListsAction() sdk.Action {
	return &GetSubscriberListsAction{}
}
