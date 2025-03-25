# JSON Converter Integration Documentation

The JSON Converter integration provides tools for transforming and processing JSON data within your workflows. This integration is especially useful for data transformation, validation, and preparing payloads for other integrations or APIs.

## Overview

JSON is the standard format for data exchange in modern applications. The JSON Converter integration offers actions to help you work with JSON data efficiently, including:

- Converting different data formats to JSON
- Converting JSON to various string formats
- Formatting and validating JSON
- Transforming JSON structures for compatibility with other systems

## Available Actions

### Convert To JSON

Converts input data to properly formatted JSON. This action can handle various input formats including existing JSON strings, key-value pairs, and more. See the [full documentation](text_to_json.md) for details.

### JSON to Text

Converts JSON data to text format with various output options. Useful for preparing data for logging, display, or for systems that require string input instead of structured JSON. See the [full documentation](json_to_text.md) for details.

## Requirements

To use the JSON Converter integration, you need:

* A Wakflo account with integration capabilities enabled
* Basic understanding of JSON data structures
* Input data in a supported format (string, key-value pairs, JSON objects, etc.)

## Common Use Cases

### Two-Way Data Transformation

Use the JSON Converter actions for bi-directional transformation:

1. Use Convert To JSON to parse and structure incoming string data
2. Process or transform the structured JSON within your workflow
3. Use JSON to String to convert the processed data back to a string format for output

### API Integration Chain

Create multi-step API integration workflows:

1. Receive data from one API in JSON format
2. Convert to string format with specific separators required by another system
3. Send the formatted string to the target system
4. Process the response by converting it back to JSON

### Data Formatting for Display

Prepare JSON data for user-friendly display:

1. Format complex JSON with proper indentation for logging or debugging
2. Convert JSON to key-value format for simplified user interfaces
3. Apply custom separators to match specific display requirements

## Implementation Example

Here's a simple example of how to use both actions in a workflow:

```go
// Import the JSON Converter actions
import (
    "github.com/wakflo/extensions/internal/integrations/jsonconverter/actions"
)

// In your setup code
func RegisterActions() {
    // Register both actions
    actions.RegisterConvertToJsonAction()
    actions.RegisterJsonToStringAction()
}

// Example workflow using both actions
func ProcessData(inputData string, outputFormat string) (string, error) {
    // Step 1: Convert string input to JSON
    convertToJsonAction := actions.NewConvertToJsonAction()
    jsonInputs := map[string]interface{}{
        "inputData": inputData,
        "format": false,
    }
    
    jsonResult, err := convertToJsonAction.Perform(ctx)
    if err != nil {
        return "", err
    }
    
    // Extract the JSON data
    jsonData := jsonResult["json"].(map[string]interface{})
    
    // Step 2: Process the data as needed
    // ... your processing logic here ...
    
    // Step 3: Convert the processed JSON back to string
    jsonToStringAction := actions.NewJsonToStringAction()
    stringInputs := map[string]interface{}{
        "inputJSON": jsonData,
        "format": true,
        "outputFormat": outputFormat,
        "keyValueSeparator": "=",
        "pairSeparator": "&",
    }
    
    stringResult, err := jsonToStringAction.Perform(ctx)
    if err != nil {
        return "", err
    }
    
    // Return the final string result
    return stringResult["result"].(string), nil
}
```

## Best Practices

1. **Data Validation**: Always validate both input and output data to ensure it matches your expected structure.

2. **Error Handling**: Implement proper error handling for both parsing and formatting operations.

3. **Format Selection**: Choose the appropriate output format based on the requirements of your target system:
   - Use JSON format for systems that expect structured data
   - Use key-value format for simpler systems or URL parameters

4. **Nested Data**: Be careful with deeply nested structures when using key-value formats; for complex structures, JSON output format is usually more appropriate.

5. **Custom Separators**: When working with key-value output, select separators that won't conflict with your data content.

## Troubleshooting

Common issues and their solutions:

1. **Formatting Errors**: If you encounter JSON formatting errors, verify that your input is valid JSON or properly formatted key-value pairs.

2. **Type Conversion Issues**: Be aware that the Convert To JSON action automatically converts numeric-looking strings to numbers, which might not be desired in all cases.

3. **Separator Conflicts**: If your data contains characters that match your separators in key-value format, they will be automatically quoted. Consider using uncommon separators if this causes issues.

4. **Missing Fields**: If fields are missing in the output, check that your JSON is properly formed and that all required fields are present.

## Support

For additional assistance with the JSON Converter integration:

* Check the detailed action documentation:
  * [Convert To JSON action documentation](./convert-to-json.md)
  * [JSON to String action documentation](./json-to-string.md)
* Contact Wakflo support at support@wakflo.com
* Visit the [Wakflo documentation](https://docs.wakflo.com) for more information