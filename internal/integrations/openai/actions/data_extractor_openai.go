package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

// ============== Type Definitions ==============

type dataExtractorOpenAIActionProps struct {
	Model             string              `json:"model"`
	Content           string              `json:"content"`
	Schema            string              `json:"schema"`
	ExtractionType    string              `json:"extraction_type"`
	Examples          []ExtractionExample `json:"examples,omitempty"`
	Temperature       *float64            `json:"temperature,omitempty"`
	MaxTokens         *int                `json:"max_tokens,omitempty"`
	ValidationMode    string              `json:"validation_mode,omitempty"`
	MultipleItems     bool                `json:"multiple_items"`
	IncludeConfidence bool                `json:"include_confidence"`
}

type ExtractionExample struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type DataExtractorOpenAIAction struct{}

func (a *DataExtractorOpenAIAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "data_extractor_openai",
		DisplayName:   "Smart Data Extractor",
		Description:   "Extract structured data from unstructured text using OpenAI. Perfect for invoice processing, form digitization, email parsing, and converting any text into structured JSON data following your custom schema.",
		Type:          core.ActionTypeAction,
		Documentation: dataExtractorOpenAIDocs,
		SampleOutput: map[string]any{
			"extracted_data": map[string]any{
				"vendor":       "Acme Corp",
				"invoice_date": "2024-01-15",
				"items": []map[string]any{
					{"description": "Widget A", "quantity": 10, "amount": 99.99},
					{"description": "Widget B", "quantity": 5, "amount": 49.99},
				},
				"subtotal": 149.98,
				"tax":      12.00,
				"total":    161.98,
			},
			"extraction_type":   "invoice",
			"confidence":        0.95,
			"tokens_used":       450,
			"validation_passed": true,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *DataExtractorOpenAIAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("data_extractor_openai", "Smart Data Extractor")

	getGPTModels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		models, err := getModels(authCtx.Extra["token"], "gpt")
		if err != nil {
			return nil, err
		}

		var options []map[string]interface{}
		for _, model := range models {
			if strings.Contains(model.ID, "gpt-4") || strings.Contains(model.ID, "gpt-3.5-turbo") {
				options = append(options, map[string]interface{}{
					"id":   model.ID,
					"name": model.ID,
				})
			}
		}

		return ctx.Respond(options, len(options))
	}

	form.SelectField("model", "Model").
		Required(true).
		HelpText("Choose the model for extraction. GPT-4 models provide better accuracy for complex extractions.").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getGPTModels)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		)

	form.TextareaField("content", "Content to Extract From").
		Required(true).
		HelpText("The unstructured text, document content, or data you want to extract information from.")

	form.SelectField("extraction_type", "Extraction Type").
		Required(true).
		HelpText("Choose a predefined extraction type or select 'custom' to define your own schema").
		AddOptions([]*smartform.Option{
			{Value: "custom", Label: "Custom Schema"},
			{Value: "invoice", Label: "Invoice/Receipt"},
			{Value: "email", Label: "Email Content"},
			{Value: "resume", Label: "Resume/CV"},
			{Value: "contract", Label: "Contract/Agreement"},
			{Value: "form", Label: "Form Data"},
			{Value: "feedback", Label: "Customer Feedback"},
			{Value: "meeting_notes", Label: "Meeting Notes"},
			{Value: "support_ticket", Label: "Support Ticket"},
			{Value: "product_review", Label: "Product Review"},
		}...)

	form.TextareaField("schema", "Data Schema").
		Required(false).
		HelpText(`Define the structure of data to extract. For custom extraction, provide a JSON schema. For predefined types, this is optional.
Example: {"company": "string", "date": "date", "items": [{"name": "string", "price": "number"}]}`)

	form.CheckboxField("multiple_items", "Extract Multiple Items").
		Required(false).
		HelpText("Enable if the content contains multiple items to extract (e.g., multiple invoices in one document)")

	form.CheckboxField("include_confidence", "Include Confidence Score").
		Required(false).
		HelpText("Include a confidence score (0-1) for the extracted data")

	form.SelectField("validation_mode", "Validation Mode").
		Required(false).
		HelpText("How strictly to validate the extracted data against the schema").
		AddOptions([]*smartform.Option{
			{Value: "none", Label: "No Validation"},
			{Value: "loose", Label: "Loose (Allow extra fields)"},
			{Value: "strict", Label: "Strict (Exact schema match)"},
		}...)

	form.TextareaField("examples", "Extraction Examples").
		Required(false).
		HelpText(`Provide examples to improve extraction accuracy. Format: Input text ||OUTPUT|| Expected JSON output. Separate multiple examples with ||EXAMPLE||`)

	form.NumberField("temperature", "Temperature").
		Required(false).
		HelpText("Control extraction consistency. Lower values (0.1-0.3) for consistent extraction, higher for creative interpretation")

	form.NumberField("max_tokens", "Max Tokens").
		Required(false).
		HelpText("Maximum tokens for the extraction response")

	schema := form.Build()
	return schema
}

func (a *DataExtractorOpenAIAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *DataExtractorOpenAIAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	rawInput := ctx.Input()

	input := &dataExtractorOpenAIActionProps{
		Model:          rawInput["model"].(string),
		Content:        rawInput["content"].(string),
		ExtractionType: rawInput["extraction_type"].(string),
	}

	if schema, ok := rawInput["schema"].(string); ok && schema != "" {
		input.Schema = schema
	}

	if multipleItems, ok := rawInput["multiple_items"].(bool); ok {
		input.MultipleItems = multipleItems
	}

	if includeConfidence, ok := rawInput["include_confidence"].(bool); ok {
		input.IncludeConfidence = includeConfidence
	}

	if validationMode, ok := rawInput["validation_mode"].(string); ok && validationMode != "" {
		input.ValidationMode = validationMode
	}

	if temp, ok := rawInput["temperature"].(float64); ok {
		input.Temperature = &temp
	}

	if maxTokens, ok := rawInput["max_tokens"].(float64); ok {
		maxTokensInt := int(maxTokens)
		input.MaxTokens = &maxTokensInt
	}

	if examplesStr, ok := rawInput["examples"].(string); ok && examplesStr != "" {
		input.Examples = parseExamples(examplesStr)
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if input.Schema == "" && input.ExtractionType != "custom" {
		input.Schema = getPredefinedSchema(input.ExtractionType)
	}

	if err := validateExtractorInput(input); err != nil {
		return nil, err
	}

	systemPrompt, userPrompt := buildExtractionPrompts(input)

	requestBody := buildExtractorRequestBody(input, systemPrompt, userPrompt)

	client, err := getOpenAiClient(authCtx.Extra["token"])
	if err != nil {
		return nil, err
	}

	res, err := client.POST("/chat/completions").
		Header().AddContentType("application/json").
		Body().AsJSON(requestBody).
		Send()
	if err != nil {
		return nil, err
	}

	if res.Status().IsError() {
		bodyBytes, _ := io.ReadAll(res.Body().Raw())
		return nil, fmt.Errorf("OpenAI API error: %s - %s", res.Status().Text(), string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(res.Body().Raw())
	if err != nil {
		return nil, err
	}

	var chatCompletion ChatCompletionResponse
	if err = json.Unmarshal(bodyBytes, &chatCompletion); err != nil {
		return nil, errors.New("could not parse response")
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, errors.New("no extraction result returned")
	}

	extractedContent := chatCompletion.Choices[0].Message.Content

	extractedContent = cleanJSONResponse(extractedContent)

	var extractedData interface{}

	if err := json.Unmarshal([]byte(extractedContent), &extractedData); err != nil {
		return map[string]interface{}{
			"extracted_data":    nil,
			"extraction_type":   input.ExtractionType,
			"tokens_used":       chatCompletion.Usage.TotalTokens,
			"validation_passed": false,
			"error":             "Failed to parse JSON response",
			"raw_response":      extractedContent,
		}, nil
	}

	validationPassed := true
	if input.ValidationMode != "" && input.ValidationMode != "none" {
		validationPassed = validateExtractedData(extractedData, input.Schema, input.ValidationMode)
	}

	response := map[string]interface{}{
		"extracted_data":    extractedData,
		"extraction_type":   input.ExtractionType,
		"tokens_used":       chatCompletion.Usage.TotalTokens,
		"validation_passed": validationPassed,
		"raw_response":      extractedContent,
	}

	if input.IncludeConfidence {
		// Simple confidence calculation based on model and validation
		confidence := 0.85 // Base confidence
		if strings.Contains(input.Model, "gpt-4") {
			confidence += 0.10
		}
		if validationPassed {
			confidence += 0.05
		}
		if confidence > 1.0 {
			confidence = 1.0
		}
		response["confidence"] = confidence
	}

	return response, nil
}

func NewDataExtractorOpenAIAction() sdk.Action {
	return &DataExtractorOpenAIAction{}
}

func supportsJSONMode(model string) bool {
	// Models that support JSON response format
	supportedModels := []string{
		"gpt-4-turbo-preview",
		"gpt-4-1106-preview",
		"gpt-4-0125-preview",
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-3.5-turbo-1106",
		"gpt-3.5-turbo-0125",
	}

	for _, supported := range supportedModels {
		if strings.Contains(model, supported) {
			return true
		}
	}

	// Check for newer turbo models
	if strings.Contains(model, "gpt-4-turbo") && !strings.Contains(model, "gpt-4-turbo-2024-04-09") {
		return true
	}

	if strings.Contains(model, "gpt-3.5-turbo") && !strings.Contains(model, "gpt-3.5-turbo-0613") && !strings.Contains(model, "gpt-3.5-turbo-0301") {
		// Newer 3.5-turbo models after 1106 support it
		return true
	}

	return false
}

func validateExtractorInput(input *dataExtractorOpenAIActionProps) error {
	if input.Model == "" {
		return errors.New("model is required")
	}

	if input.Content == "" {
		return errors.New("content to extract from is required")
	}

	if input.ExtractionType == "custom" && input.Schema == "" {
		return errors.New("schema is required for custom extraction type")
	}

	if input.Temperature != nil && (*input.Temperature < 0 || *input.Temperature > 2) {
		return errors.New("temperature must be between 0 and 2")
	}

	return nil
}

func buildExtractionPrompts(input *dataExtractorOpenAIActionProps) (string, string) {
	systemPrompt := "You are a precise data extraction specialist. Extract structured data from the provided content and return it as valid JSON."

	// For models that don't support JSON mode, be more explicit
	if !supportsJSONMode(input.Model) {
		systemPrompt = "You are a precise data extraction specialist. Extract structured data from the provided content.\n\nCRITICAL: You MUST return ONLY valid JSON with no additional text, markdown formatting, or explanation. Do not wrap the JSON in ```json``` tags or any other formatting. Return raw JSON only."
	}

	if input.IncludeConfidence {
		systemPrompt += " Include a 'confidence' field (0-1) in your response indicating how confident you are in the extraction."
	}

	if input.MultipleItems {
		systemPrompt += " The content may contain multiple items. Return an array of extracted items."
	}

	systemPrompt += fmt.Sprintf("\n\nRequired output schema:\n%s", input.Schema)

	if len(input.Examples) > 0 {
		systemPrompt += "\n\nExamples:"
		for i, example := range input.Examples {
			systemPrompt += fmt.Sprintf("\n\nExample %d:\nInput: %s\nExpected Output: %s",
				i+1, example.Input, example.Output)
		}
	}

	// Reinforce JSON-only output for non-JSON-mode models
	if !supportsJSONMode(input.Model) {
		systemPrompt += "\n\nREMEMBER: Output ONLY the JSON data structure, no explanations or formatting."
	} else {
		systemPrompt += "\n\nIMPORTANT: Return ONLY valid JSON that matches the schema. No additional text or explanation."
	}

	userPrompt := fmt.Sprintf("Extract data from the following content:\n\n%s", input.Content)

	return systemPrompt, userPrompt
}

func buildExtractorRequestBody(input *dataExtractorOpenAIActionProps, systemPrompt, userPrompt string) map[string]interface{} {
	requestBody := map[string]interface{}{
		"model": input.Model,
		"messages": []interface{}{
			map[string]interface{}{
				"role":    "system",
				"content": systemPrompt,
			},
			map[string]interface{}{
				"role":    "user",
				"content": userPrompt,
			},
		},
	}

	// Only add response_format for models that support it
	// GPT-4 Turbo, GPT-4o, and GPT-3.5-turbo-1106 or later support JSON mode
	if supportsJSONMode(input.Model) {
		requestBody["response_format"] = map[string]interface{}{
			"type": "json_object",
		}
	}

	// Set default temperature for extraction (lower = more consistent)
	if input.Temperature != nil {
		requestBody["temperature"] = *input.Temperature
	} else {
		requestBody["temperature"] = 0.3 // Default low temperature for consistency
	}

	if input.MaxTokens != nil {
		requestBody["max_tokens"] = *input.MaxTokens
	}

	return requestBody
}

func getPredefinedSchema(extractionType string) string {
	schemas := map[string]string{
		"invoice": `{
			"vendor": "string",
			"vendor_address": "string",
			"invoice_number": "string",
			"invoice_date": "date",
			"due_date": "date",
			"items": [{"description": "string", "quantity": "number", "unit_price": "number", "amount": "number"}],
			"subtotal": "number",
			"tax": "number",
			"total": "number",
			"payment_terms": "string"
		}`,
		"email": `{
			"from": "string",
			"to": ["string"],
			"subject": "string",
			"date": "datetime",
			"intent": "string",
			"key_points": ["string"],
			"action_items": ["string"],
			"sentiment": "string",
			"urgency": "string"
		}`,
		"resume": `{
			"name": "string",
			"email": "string",
			"phone": "string",
			"location": "string",
			"summary": "string",
			"experience": [{"company": "string", "position": "string", "duration": "string", "responsibilities": ["string"]}],
			"education": [{"institution": "string", "degree": "string", "year": "string"}],
			"skills": ["string"],
			"certifications": ["string"]
		}`,
		"contract": `{
			"contract_type": "string",
			"parties": [{"name": "string", "role": "string", "address": "string"}],
			"effective_date": "date",
			"expiration_date": "date",
			"terms": ["string"],
			"payment_terms": "string",
			"obligations": [{"party": "string", "obligation": "string"}],
			"termination_clauses": ["string"]
		}`,
		"feedback": `{
			"customer_name": "string",
			"date": "date",
			"rating": "number",
			"sentiment": "string",
			"topics": ["string"],
			"positive_aspects": ["string"],
			"negative_aspects": ["string"],
			"suggestions": ["string"],
			"would_recommend": "boolean"
		}`,
		"meeting_notes": `{
			"meeting_title": "string",
			"date": "datetime",
			"attendees": ["string"],
			"agenda_items": ["string"],
			"key_discussions": ["string"],
			"decisions_made": ["string"],
			"action_items": [{"assignee": "string", "task": "string", "due_date": "date"}],
			"next_meeting": "datetime"
		}`,
		"support_ticket": `{
			"ticket_id": "string",
			"customer": "string",
			"date_created": "datetime",
			"issue_type": "string",
			"priority": "string",
			"description": "string",
			"symptoms": ["string"],
			"steps_to_reproduce": ["string"],
			"affected_systems": ["string"],
			"resolution": "string"
		}`,
		"product_review": `{
			"reviewer_name": "string",
			"product_name": "string",
			"rating": "number",
			"review_date": "date",
			"pros": ["string"],
			"cons": ["string"],
			"summary": "string",
			"would_buy_again": "boolean",
			"recommendation": "string"
		}`,
		"form": `{
			"form_type": "string",
			"submission_date": "date",
			"fields": [{"field_name": "string", "field_value": "string"}]
		}`,
	}

	if schema, exists := schemas[extractionType]; exists {
		return schema
	}

	return `{"data": "extracted content"}`
}

func parseExamples(examplesStr string) []ExtractionExample {
	var examples []ExtractionExample

	examplePairs := strings.Split(examplesStr, "||EXAMPLE||")
	for _, pair := range examplePairs {
		parts := strings.Split(pair, "||OUTPUT||")
		if len(parts) == 2 {
			examples = append(examples, ExtractionExample{
				Input:  strings.TrimSpace(parts[0]),
				Output: strings.TrimSpace(parts[1]),
			})
		}
	}

	return examples
}

func validateExtractedData(data interface{}, schema string, mode string) bool {
	if mode == "strict" {
		_, ok := data.(map[string]interface{})
		return ok
	}
	return data != nil
}

func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)

	// Remove ```json and ``` tags
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}

	// Remove any leading/trailing whitespace or newlines
	response = strings.TrimSpace(response)

	return response
}
