package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type GetSubscriberListsAction struct{}

// Metadata returns metadata about the action
func (a *GetSubscriberListsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_subscriber_lists",
		DisplayName:   "Get Subscriber Lists",
		Description:   "Retrieve all subscriber lists for a client.",
		Type:          core.ActionTypeAction,
		Documentation: getSubscriberList,
		Icon:          "mdi:format-list-bulleted",
		SampleOutput: map[string]interface{}{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetSubscriberListsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_subscriber_lists", "Get Subscriber Lists")

	form.NumberField("page", "Page").
		Placeholder("Enter page number").
		Required(false).
		HelpText("The page number to retrieve (for pagination).")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetSubscriberListsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetSubscriberListsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	clientID, clientIDExists := authCtx.Extra["client-id"]
	apiKey, apiKeyExists := authCtx.Extra["api-key"]

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

func NewGetSubscriberListsAction() sdk.Action {
	return &GetSubscriberListsAction{}
}
