package actions

import (
	"encoding/json"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type convertToJSONActionProps struct {
	InputData string `json:"inputData"`
	Format    bool   `json:"format"`
}

type ConvertToJSONAction struct{}

func (a *ConvertToJSONAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "convert_to_json",
		DisplayName:   "Convert To JSON",
		Description:   "Converts input data to properly formatted JSON. Useful for data integrations and transformations.",
		Type:          core.ActionTypeAction,
		Documentation: textToJSONDocs,
		Icon:          "json",
		SampleOutput: map[string]any{
			"json": map[string]interface{}{
				"name":  "John Doe",
				"age":   30,
				"email": "john@example.com",
			},
			"type": "key-value",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ConvertToJSONAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("convert_to_json", "Convert To JSON")

	form.TextareaField("inputData", "Input Data").
		Required(true).
		HelpText("The data to convert to JSON. Can be a JSON string, CSV data, or key-value pairs.")

	form.CheckboxField("format", "Pretty Print").
		Required(false).
		HelpText("Format the JSON output with indentation for better readability.")

	schema := form.Build()

	return schema
}

func (a *ConvertToJSONAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *ConvertToJSONAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[convertToJSONActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.InputData == "" {
		return nil, fmt.Errorf("input data cannot be empty")
	}

	// First try to parse as JSON
	var parsedData interface{}
	err = json.Unmarshal([]byte(input.InputData), &parsedData)

	if err == nil {
		if input.Format {
			formattedJSON, err := json.MarshalIndent(parsedData, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("error formatting JSON: %v", err)
			}
			return map[string]interface{}{
				"json": string(formattedJSON),
				"type": "json",
			}, nil
		}

		return map[string]interface{}{
			"json": parsedData,
			"type": "json",
		}, nil
	}

	result, err := parseKeyValuePairs(input.InputData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse input data: %v", err)
	}

	if input.Format {
		formattedJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error formatting JSON: %v", err)
		}
		return map[string]interface{}{
			"json": string(formattedJSON),
			"type": "key-value",
		}, nil
	}

	return map[string]interface{}{
		"json": result,
		"type": "key-value",
	}, nil
}

// parseKeyValuePairs parses input like "key1:value1,key2:value2" into a map
func parseKeyValuePairs(input string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Split by commas to get key-value pairs
	pairs := splitOnCommas(input)

	for _, pair := range pairs {
		// Split each pair by colon
		kv := splitOnColons(pair)
		if len(kv) != 2 {
			continue // Skip invalid pairs
		}

		key := kv[0]
		value := kv[1]

		// Try to convert value to number if possible
		if num, err := parseNumber(value); err == nil {
			result[key] = num
		} else {
			// Otherwise keep as string
			result[key] = value
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid key-value pairs found")
	}

	return result, nil
}

// splitOnCommas splits a string by commas, but respects quoted sections
func splitOnCommas(input string) []string {
	var result []string
	var current string
	inQuotes := false

	for _, char := range input {
		if char == '"' {
			inQuotes = !inQuotes
			current += string(char)
		} else if char == ',' && !inQuotes {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

// splitOnColons splits a string by the first colon
func splitOnColons(input string) []string {
	for i, char := range input {
		if char == ':' {
			return []string{
				input[:i],
				input[i+1:],
			}
		}
	}
	return []string{input}
}

// parseNumber tries to parse a string as a number
func parseNumber(s string) (interface{}, error) {
	var i int
	if _, err := fmt.Sscanf(s, "%d", &i); err == nil {
		return i, nil
	}

	var f float64
	if _, err := fmt.Sscanf(s, "%f", &f); err == nil {
		return f, nil
	}

	return nil, fmt.Errorf("not a number")
}

func NewConvertToJSONAction() sdk.Action {
	return &ConvertToJSONAction{}
}
