package shared

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("activecampaign-auth", "ActiveCampaign API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api_url", "API URL (Required)").
		Required(true).
		HelpText("Your ActiveCampaign API URL (e.g., https://youraccountname.api-us1.com)")

	_ = form.TextField("api_key", "API Key (Required)").
		Required(true).
		HelpText("Your ActiveCampaign API Key.")

	ActiveCampaignSharedAuth = form.Build()
)

func RegisterActiveCampaignListsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getLists := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiURL, apiURLExists := authCtx.Extra["api_url"]
		apiKey, apiKeyExists := authCtx.Extra["api_key"]
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

	return form.SelectField("list-id", "ActiveCampaign List").
		Placeholder("Select an ActiveCampaign list").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getLists)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select an ActiveCampaign list")
}

func RegisterContactsProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getContacts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiURL, apiURLExists := authCtx.Extra["api_url"]
		apiKey, apiKeyExists := authCtx.Extra["api_key"]
		if !apiURLExists || !apiKeyExists || apiURL == "" || apiKey == "" {
			return nil, errors.New("API URL and API Key are required")
		}

		result, err := GetActiveCampaignClient(apiURL, apiKey, "contacts")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch contacts: %v", err)
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected result format from API")
		}

		contactsData, ok := resultMap["contacts"]
		if !ok {
			return nil, errors.New("no contacts found in response")
		}

		contactItems, ok := contactsData.([]interface{})
		if !ok {
			return nil, errors.New("unexpected format for contacts data")
		}

		items := make([]map[string]any, 0, len(contactItems))
		for _, item := range contactItems {
			contact, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			var id string
			var email string

			// Extract ID
			if idStr, ok := contact["id"].(string); ok {
				id = idStr
			} else if idNum, ok := contact["id"].(float64); ok {
				id = strconv.Itoa(int(idNum))
			}

			// Extract email
			if emailStr, ok := contact["email"].(string); ok {
				email = emailStr
			}

			if id != "" && email != "" {
				items = append(items, map[string]any{
					"id":   id,
					"name": email,
				})
			}
		}

		return ctx.Respond(items, len(items))
	}

	return form.SelectField("contact-id", "Contact").
		Placeholder("Select contact").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getContacts)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a contact")
}
