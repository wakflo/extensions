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
	form = smartform.NewAuthForm("convertKit-auth", "ConvertKit API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "ConvertKit API Key (Required)").
		Required(true).
		HelpText("ConvertKit API Key")

	_ = form.TextField("api-secret", "API Secret (Required)").
		Required(true).
		HelpText("ConvertKit API Secret")

	ConvertKitSharedAuth = form.Build()
)

func RegisterTagsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getTags := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Get API key from auth context
		apiKey, ok := authCtx.Extra["api-key"]
		if !ok {
			return nil, errors.New("API key not found in auth context")
		}

		// ConvertKit API URL with query parameter
		apiURL := fmt.Sprintf("https://api.convertkit.com/v3/tags?api_key=%s", apiKey)

		// Create HTTP request
		req, err := http.NewRequest(http.MethodGet, apiURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Define response structure based on ConvertKit API
		var response struct {
			Tags []struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				CreatedAt string `json:"created_at"`
			} `json:"tags"`
		}

		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}

		var options []map[string]interface{}
		for _, tag := range response.Tags {
			options = append(options, map[string]interface{}{
				"id":   fmt.Sprintf("%d", tag.ID),
				"name": tag.Name,
			})
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField("tag_id", "Tag").
		Placeholder("Select a tag").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTags)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select a ConvertKit tag to apply")
}

func RegisterSubscribersProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getSubscribers := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Get API secret key from auth context
		apiSecret, ok := authCtx.Extra["api-secret"]
		if !ok {
			return nil, errors.New("API secret key not found in auth context")
		}

		// ConvertKit API URL with query parameter
		apiURL := fmt.Sprintf("https://api.convertkit.com/v3/subscribers?api_secret=%s", apiSecret)

		// Create HTTP request
		req, err := http.NewRequest(http.MethodGet, apiURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Define response structure based on ConvertKit API
		var response SubscribersResponse

		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}

		var options []map[string]interface{}
		for _, subscriber := range response.Subscribers {
			// Create display name with email and name (if available)
			displayName := subscriber.EmailAddress
			if subscriber.FirstName != "" {
				displayName = fmt.Sprintf("%s (%s)", subscriber.FirstName, subscriber.EmailAddress)
			}

			options = append(options, map[string]interface{}{
				"id":   fmt.Sprintf("%d", subscriber.ID),
				"name": displayName,
			})
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField("subscriber_id", "Subscriber").
		Placeholder("Select a subscriber").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSubscribers)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select a ConvertKit subscriber")
}
