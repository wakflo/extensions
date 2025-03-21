# Convert To JSON

## Description

Converts input data to properly formatted JSON. Useful for data integrations and transformations. The action can handle valid JSON strings, CSV data, and key-value pairs.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `inputData` | String | Yes | The data to convert to JSON. Can be a JSON string, CSV data, or key-value pairs. |
| `format` | Boolean | No | Format the JSON output with indentation for better readability. |

## Details

- **Type**: sdkcore.ActionTypeNormal
- **Icon**: json

## Output

Returns a JSON object containing the parsed data and the type of parsing that was performed (json or key-value).

## Sample Output

```json
{
  "json": {
    "name": "John Doe",
    "age": 30,
    "email": "john@example.com"
  },
  "type": "key-value"
}
```

When the `format` option is enabled and key-value input is provided:

```json
{
  "json": "{\n  \"name\": \"John Doe\",\n  \"age\": 30,\n  \"email\": \"john@example.com\"\n}",
  "type": "key-value"
}
```

## Example Usage

### Converting key-value pairs to JSON

**Input:**
```
name:John Doe,age:30,email:john@example.com
```

**Output:**
```json
{
  "json": {
    "name": "John Doe",
    "age": 30,
    "email": "john@example.com"
  },
  "type": "key-value"
}
```

### Formatting existing JSON

**Input:**
```json
{"name":"John Doe","age":30,"email":"john@example.com"}
```

With `format` set to `true`, the output will be:

```json
{
  "json": "{\n  \"name\": \"John Doe\",\n  \"age\": 30,\n  \"email\": \"john@example.com\"\n}",
  "type": "json"
}
```

## Notes

- The action automatically detects numeric values and converts them to the appropriate type
- When passing already valid JSON, the action validates and optionally reformats it
- For key-value pairs, use the format `key1:value1,key2:value2`
- The action handles quoted sections in key-value pairs, so values containing commas can be enclosed in quotes