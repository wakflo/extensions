package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("campaignmonitor-auth", "Campaign Monitor API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "API Key (Required)").
		Required(true).
		HelpText("Your Campaign Monitor API key")

	_ = form.TextField("client-id", "Client ID (Required)").
		Required(true).
		HelpText("Your Campaign Monitor Client ID")

	CampaignMonitorSharedAuth = form.Build()
)

func RegisterCreateSendSubscriberListsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getLists := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		clientID, clientIDExists := authCtx.Extra["client-id"]
		apiKey, apiKeyExists := authCtx.Extra["api-key"]

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

	return form.SelectField("subscriber_list", "Subscriber Lists").
		Placeholder("Select a subscriber list").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getLists)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select CreateSend subscriber list")
}

type EmailItem struct {
	Value string `json:"value"`
}
