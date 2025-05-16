package shared

import (
	"encoding/json"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

func GetProductProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	listSendOwlProducts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get API credentials from context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiKey := authCtx.Extra["api_key"]
		apiSecret := authCtx.Extra["api_secret"]

		// SendOwl API endpoint for fetching products
		endpoint := "/products"

		// Use the SendOwl client to fetch products
		response, err := GetSendOwlClient(BaseURL, apiKey, apiSecret, endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch SendOwl products: %v", err)
		}

		// Process the response
		var options []map[string]interface{}

		if response.IsArray {
			// The response structure is an array of objects with a "product" key
			for _, item := range response.Array {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				// Extract the product map from the item
				productRaw, ok := itemMap["product"].(map[string]interface{})
				if !ok {
					continue
				}

				// Extract product properties
				id, idOk := productRaw["id"]
				name, nameOk := productRaw["name"].(string)

				if !idOk || !nameOk {
					continue
				}

				// Convert ID to string based on type
				var idStr string
				switch v := id.(type) {
				case float64:
					idStr = fmt.Sprintf("%.0f", v)
				case string:
					idStr = v
				case json.Number:
					idStr = string(v)
				default:
					idStr = fmt.Sprintf("%v", v)
				}

				options = append(options, map[string]interface{}{
					"value": idStr,
					"label": name,
				})

			}
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select product").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&listSendOwlProducts)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}

func GetOrderProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	listSendOwlProducts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get API credentials from context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		apiKey := authCtx.Extra["api_key"]
		apiSecret := authCtx.Extra["api_secret"]

		// SendOwl API endpoint for fetching products
		endpoint := "/orders"

		// Use the SendOwl client to fetch products
		response, err := GetSendOwlClient(AltBaseURL, apiKey, apiSecret, endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch SendOwl products: %v", err)
		}

		// Process the response
		var options []map[string]interface{}

		if response.IsArray {
			// The response structure is an array of objects with a "order" key
			for _, item := range response.Array {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				// Extract the order map from the item
				orderRaw, ok := itemMap["order"].(map[string]interface{})
				if !ok {
					continue
				}

				// Extract product properties
				id, idOk := orderRaw["id"]
				if !idOk {
					continue
				}

				// Convert ID to string based on type
				var idStr string
				switch v := id.(type) {
				case float64:
					idStr = fmt.Sprintf("%.0f", v)
				case string:
					idStr = v
				case json.Number:
					idStr = string(v)
				default:
					idStr = fmt.Sprintf("%v", v)
				}

				options = append(options, map[string]interface{}{
					"value": idStr,
					"label": "#SO" + idStr,
				})

			}
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select an order").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&listSendOwlProducts)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}
