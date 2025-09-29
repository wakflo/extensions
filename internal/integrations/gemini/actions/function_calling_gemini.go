package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type functionCallingActionProps struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	FunctionSchema string `json:"function_schema"`
	SystemPrompt   string `json:"system_prompt"`
}

type FunctionCallingAction struct{}

func (a *FunctionCallingAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "function_calling_gemini",
		DisplayName:   "Extract Structured Data",
		Description:   "Extract structured data from text using Gemini's function calling capabilities. Perfect for parsing emails, documents, or any unstructured text into JSON.",
		Type:          core.ActionTypeAction,
		Documentation: functionCallingGeminiDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"extracted_data": map[string]any{
				"name":  "John Doe",
				"email": "john@example.com",
				"phone": "+1234567890",
			},
			"model": "gemini-1.5-flash",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *FunctionCallingAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("function_calling_gemini", "Extract Structured Data")

	form.TextareaField("prompt", "Input Text").
		Placeholder("Paste the text you want to extract data from").
		HelpText("The text containing the information to extract").
		Required(true)

	form.TextareaField("function_schema", "Extraction Schema").
		Placeholder(`{
  "name": {"type": "string", "description": "Person's name"},
  "email": {"type": "string", "description": "Email address"},
  "amount": {"type": "number", "description": "Dollar amount"}
}`).
		HelpText("JSON schema defining the fields to extract (each field should have type and description)").
		Required(true)

	form.TextareaField("system_prompt", "Extraction Instructions").
		Placeholder("Extract the customer information from this email").
		HelpText("Additional instructions for the extraction (optional)").
		Required(false)

	RegisterModelProps(form)

	return form.Build()
}

func (a *FunctionCallingAction) Auth() *core.AuthMetadata {
	return nil
}

// convertJSONSchemaToGenaiSchema converts a JSON schema to Gemini Schema format
func convertJSONSchemaToGenaiSchema(jsonSchema map[string]interface{}) (map[string]*genai.Schema, error) {
	properties := make(map[string]*genai.Schema)

	for key, value := range jsonSchema {
		propMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid property format for field '%s'", key)
		}

		schema := &genai.Schema{}

		// Get type
		if typeStr, ok := propMap["type"].(string); ok {
			switch typeStr {
			case "string":
				schema.Type = genai.TypeString
			case "number", "integer":
				schema.Type = genai.TypeNumber
			case "boolean":
				schema.Type = genai.TypeBoolean
			case "array":
				schema.Type = genai.TypeArray
				// If array has items definition
				if items, ok := propMap["items"].(map[string]interface{}); ok {
					if itemType, ok := items["type"].(string); ok {
						itemSchema := &genai.Schema{}
						switch itemType {
						case "string":
							itemSchema.Type = genai.TypeString
						case "number":
							itemSchema.Type = genai.TypeNumber
						case "boolean":
							itemSchema.Type = genai.TypeBoolean
						}
						schema.Items = itemSchema
					}
				}
			case "object":
				schema.Type = genai.TypeObject
			default:
				schema.Type = genai.TypeString // default to string
			}
		} else {
			schema.Type = genai.TypeString // default to string if no type specified
		}

		// Get description
		if desc, ok := propMap["description"].(string); ok {
			schema.Description = desc
		}

		// Check if required
		if required, ok := propMap["required"].(bool); ok {
			if required {
				schema.Required = []string{key}
			}
		}

		properties[key] = schema
	}

	return properties, nil
}

func (a *FunctionCallingAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[functionCallingActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Parse the function schema
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(input.FunctionSchema), &schema); err != nil {
		return nil, fmt.Errorf("invalid function schema JSON: %w", err)
	}

	// Convert JSON schema to Gemini schema format
	properties, err := convertJSONSchemaToGenaiSchema(schema)
	if err != nil {
		return nil, err
	}

	gctx := context.Background()
	client, err := CreateGeminiClient(gctx, authCtx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Clean model name
	modelName := strings.TrimPrefix(input.Model, "models/")

	// Create function declaration with properly typed schema
	extractFunction := &genai.FunctionDeclaration{
		Name:        "extract_data",
		Description: "Extract structured data from the provided text",
		Parameters: &genai.Schema{
			Type:       genai.TypeObject,
			Properties: properties,
		},
	}

	// Configure model with function
	model := client.GenerativeModel(modelName)
	model.Tools = []*genai.Tool{
		{FunctionDeclarations: []*genai.FunctionDeclaration{extractFunction}},
	}

	// Force the model to use the function
	model.ToolConfig = &genai.ToolConfig{
		FunctionCallingConfig: &genai.FunctionCallingConfig{
			Mode:                 genai.FunctionCallingAny,
			AllowedFunctionNames: []string{"extract_data"},
		},
	}

	// Build the prompt
	fullPrompt := input.Prompt
	if input.SystemPrompt != "" {
		fullPrompt = fmt.Sprintf("%s\n\nText to analyze:\n%s", input.SystemPrompt, input.Prompt)
	}

	// Generate content with function calling
	resp, err := model.GenerateContent(gctx, genai.Text(fullPrompt))
	if err != nil {
		return nil, fmt.Errorf("failed to extract data: %w", err)
	}

	// Extract function call results
	var extractedData map[string]interface{}
	if resp != nil && len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if funcCall, ok := part.(*genai.FunctionCall); ok {
				if funcCall.Name == "extract_data" {
					extractedData = funcCall.Args
					break
				}
			}
		}
	}

	// If no function call was made, try to extract from text response
	if extractedData == nil {
		var textResponse string
		if resp != nil && len(resp.Candidates) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				if text, ok := part.(genai.Text); ok {
					textResponse += string(text)
				}
			}
		}

		// Try to parse as JSON if we got a text response
		if textResponse != "" {
			// Sometimes the model returns JSON in the text instead of using function calling
			json.Unmarshal([]byte(textResponse), &extractedData)
		}
	}

	// Return the results
	result := map[string]interface{}{
		"model": modelName,
	}

	if extractedData != nil {
		result["extracted_data"] = extractedData
		result["success"] = true
	} else {
		result["extracted_data"] = nil
		result["success"] = false
		result["message"] = "Failed to extract structured data. Try adjusting the schema or prompt."
	}

	return result, nil
}

func NewFunctionCallingAction() sdk.Action {
	return &FunctionCallingAction{}
}
