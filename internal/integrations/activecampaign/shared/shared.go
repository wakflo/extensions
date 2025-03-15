package shared

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	// "github.com/wakflo/go-sdk/sdk"
)

var SharedAuth = autoform.NewAuth().NewCustomAuth().
	SetDescription("ActiveCampaign API Authentication").
	SetLabel("ActiveCampaign Authentication").
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api_url": autoform.NewShortTextField().
			SetDisplayName("API URL").
			SetDescription("Your ActiveCampaign API URL (e.g., https://youraccountname.api-us1.com)").
			SetRequired(true).
			Build(),
		"api_key": autoform.NewShortTextField().
			SetDisplayName("API Key").
			SetDescription("Your ActiveCampaign API Key").
			SetRequired(true).
			Build(),
	}).
	Build()

func GetActiveCampaignListsInput() *sdkcore.AutoFormSchema {
	getLists := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		apiURL, apiURLExists := ctx.Auth.Extra["api_url"]
		apiKey, apiKeyExists := ctx.Auth.Extra["api_key"]

		if !apiURLExists || !apiKeyExists || apiURL == "" || apiKey == "" {
			return nil, errors.New("API URL and API Key are required")
		}

		result, err := GetActiveCampaignClient(apiURL, apiKey, "lists")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch lists: %v", err)
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected result format from API")
		}

		listsData, ok := resultMap["lists"]
		if !ok {
			return nil, errors.New("no lists found in response")
		}

		listItems, ok := listsData.([]interface{})
		if !ok {
			return nil, errors.New("unexpected format for lists data")
		}
		items := make([]map[string]any, 0, len(listItems))
		for _, item := range listItems {
			list, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			var id string
			var name string

			if idStr, ok := list["id"].(string); ok {
				id = idStr
			} else if idNum, ok := list["id"].(float64); ok {
				id = strconv.Itoa(int(idNum))
			}

			if nameStr, ok := list["name"].(string); ok {
				name = nameStr
			}

			if id != "" && name != "" {
				items = append(items, map[string]any{
					"id":   id,
					"name": name,
				})
			}
		}

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Lists").
		SetDescription("Select ActiveCampaign list").
		SetDynamicOptions(&getLists).
		SetRequired(false).Build()
}
