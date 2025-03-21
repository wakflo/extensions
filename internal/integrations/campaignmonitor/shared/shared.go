package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var SharedAuth = autoform.NewAuth().NewCustomAuth().
	SetDescription("ActiveCampaign API Authentication").
	SetLabel("ActiveCampaign Authentication").
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().
			SetDisplayName("API Key").
			SetDescription("Your Campaign Monitor API key").
			SetRequired(true).
			Build(),
		"client-id": autoform.NewShortTextField().
			SetDisplayName("Client ID").
			SetDescription("Your Campaign Monitor Client ID").
			SetRequired(true).
			Build(),
	}).
	Build()

func GetCreateSendSubscriberListsInput() *sdkcore.AutoFormSchema {
	getLists := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		clientID, clientIDExists := ctx.Auth.Extra["client-id"]
		apiKey, apiKeyExists := ctx.Auth.Extra["api-key"]

		if !clientIDExists || !apiKeyExists || clientID == "" || apiKey == "" {
			return nil, errors.New("client ID and API Key are required")
		}
		apiURL := fmt.Sprintf("https://api.createsend.com/api/v3.3/clients/%s/lists.json", clientID)

		req, err := http.NewRequest(http.MethodGet, apiURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		req.SetBasicAuth(apiKey, "")

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

		items := make([]map[string]any, 0, len(lists))
		for _, list := range lists {
			id, idOk := list["ListID"].(string)
			name, nameOk := list["Name"].(string)

			if idOk && nameOk && id != "" && name != "" {
				items = append(items, map[string]any{
					"id":   id,
					"name": name,
				})
			}
		}

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Subscriber Lists").
		SetDescription("Select CreateSend subscriber list").
		SetDynamicOptions(&getLists).
		SetRequired(false).Build()
}

type EmailItem struct {
	Value string `json:"value"` // The actual field name may be different
}
