package actions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type jsonToStringActionProps struct {
	InputJSON         interface{} `json:"inputJSON"`
	Format            bool        `json:"format"`
	OutputFormat      string      `json:"outputFormat"`
	KeyValueSeparator string      `json:"keyValueSeparator"`
	PairSeparator     string      `json:"pairSeparator"`
}

type JsonToStringAction struct{}

func (a *JsonToStringAction) Name() string {
	return "JSON to Text"
}

func (a *JsonToStringAction) Description() string {
	return "Converts JSON data to text format. Options for outputting as a JSON string, key-value pairs, or other formats."
}

func (a *JsonToStringAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *JsonToStringAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &jsonToTextDocs,
	}
}

func (a *JsonToStringAction) Icon() *string {
	icon := "json"
	return &icon
}

func (a *JsonToStringAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"inputJSON": autoform.NewCodeEditorField().
			SetDisplayName("Input JSON").
			SetDescription("The JSON data to convert to a string.").
			SetRequired(true).
			Build(),
		"format": autoform.NewBooleanField().
			SetDisplayName("Pretty Print").
			SetDescription("Format the JSON output with indentation for better readability. Only applies to JSON output format.").
			SetRequired(false).
			Build(),
		"outputFormat": autoform.NewSelectField().
			SetDisplayName("Output Format").
			SetDescription("The format of the output string.").
			SetOptions([]*sdkcore.AutoFormSchema{
				{Title: "JSON", Const: "json"},
				{Title: "Key-Value Pairs", Const: "keyvalue"},
			}).
			SetDefaultValue("json").
			SetRequired(true).
			Build(),
		"keyValueSeparator": autoform.NewShortTextField().
			SetDisplayName("Key-Value Separator").
			SetDescription("Character(s) to use between keys and values. Only applies to key-value output format.").
			SetDefaultValue(":").
			SetRequired(false).
			Build(),
		"pairSeparator": autoform.NewShortTextField().
			SetDisplayName("Pair Separator").
			SetDescription("Character(s) to use between key-value pairs. Only applies to key-value output format.").
			SetDefaultValue(",").
			SetRequired(false).
			Build(),
	}
}

func (a *JsonToStringAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[jsonToStringActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Convert inputJSON to a map[string]interface{} regardless of how it was passed
	var jsonData map[string]interface{}

	switch v := input.InputJSON.(type) {
	case map[string]interface{}:
		// Already in the right format
		jsonData = v
	case string:
		// Parse the JSON string
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid JSON string: %v", err)
		}
	default:
		// Try to marshal and unmarshal as a last resort
		bytes, err := json.Marshal(input.InputJSON)
		if err != nil {
			return nil, fmt.Errorf("unable to process JSON input: %v", err)
		}
		if err := json.Unmarshal(bytes, &jsonData); err != nil {
			return nil, fmt.Errorf("unable to process JSON input: %v", err)
		}
	}

	if len(jsonData) == 0 {
		return nil, fmt.Errorf("input JSON cannot be empty")
	}

	var result string

	switch input.OutputFormat {
	case "json":
		// Convert to JSON string
		var jsonBytes []byte
		var err error

		if input.Format {
			jsonBytes, err = json.MarshalIndent(jsonData, "", "  ")
		} else {
			jsonBytes, err = json.Marshal(jsonData)
		}

		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON: %v", err)
		}
		result = string(jsonBytes)

	case "keyvalue":
		// Convert to key-value pairs
		kvSeparator := input.KeyValueSeparator
		if kvSeparator == "" {
			kvSeparator = ":"
		}

		pairSeparator := input.PairSeparator
		if pairSeparator == "" {
			pairSeparator = ","
		}

		result = mapToKeyValueString(jsonData, kvSeparator, pairSeparator)

	default:
		return nil, fmt.Errorf("unsupported output format: %s", input.OutputFormat)
	}

	return map[string]interface{}{
		"result": result,
		"format": input.OutputFormat,
	}, nil
}

// mapToKeyValueString converts a map to a key-value string with the specified separators
func mapToKeyValueString(m map[string]interface{}, kvSeparator, pairSeparator string) string {
	pairs := make([]string, 0, len(m))

	for k, v := range m {
		// Handle different value types
		var valueStr string
		switch val := v.(type) {
		case string:
			// If the string contains the pair separator, we should quote it
			if strings.Contains(val, pairSeparator) {
				valueStr = fmt.Sprintf("\"%s\"", val)
			} else {
				valueStr = val
			}
		case map[string]interface{}:
			// For nested objects, convert to JSON
			nestedJSON, _ := json.Marshal(val)
			valueStr = string(nestedJSON)
		case []interface{}:
			// For arrays, convert to JSON
			arrayJSON, _ := json.Marshal(val)
			valueStr = string(arrayJSON)
		default:
			// For other types, use simple string conversion
			valueStr = fmt.Sprintf("%v", val)
		}

		pairs = append(pairs, fmt.Sprintf("%s%s%s", k, kvSeparator, valueStr))
	}

	return strings.Join(pairs, pairSeparator)
}

func (a *JsonToStringAction) Auth() *sdk.Auth {
	return nil
}

func (a *JsonToStringAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"result": "name:John Doe,age:30,email:john@example.com",
		"format": "keyvalue",
	}
}

func (a *JsonToStringAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewJSONToStringAction() sdk.Action {
	return &JsonToStringAction{}
}
