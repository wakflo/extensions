package actions

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

// ============== Type Definitions ==============

type visionOpenAIActionProps struct {
	Model          string   `json:"model"`
	ImageInput     string   `json:"image_input"`
	ImageType      string   `json:"image_type"`
	Prompt         string   `json:"prompt"`
	SystemPrompt   string   `json:"system_prompt,omitempty"`
	Detail         string   `json:"detail,omitempty"`
	MaxTokens      *int     `json:"max_tokens,omitempty"`
	Temperature    *float64 `json:"temperature,omitempty"`
	MultipleImages []string `json:"multiple_images,omitempty"`
}

type VisionOpenAIAction struct{}

func (a *VisionOpenAIAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "vision_openai",
		DisplayName:   "Analyze Image (Vision)",
		Description:   "Analyze images using OpenAI's GPT-4 Vision to extract text, describe content, answer questions about images, or compare multiple images. Perfect for OCR, document processing, quality control, and visual content moderation.",
		Type:          core.ActionTypeAction,
		Documentation: VisionOpenAIDocs,
		SampleOutput: map[string]any{
			"vision_analysis": "The image shows an invoice from Acme Corp dated January 20, 2024. Invoice number INV-2024-0892. Total amount is $483.84 including tax. Items listed are Premium Plan ($299), API Addon ($99), and Support Package ($50).",
			"model":           "gpt-4-vision-preview",
			"prompt":          "Extract all text and key information from this invoice",
			"image_type":      "url",
			"detail_level":    "high",
			"tokens_used":     650,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *VisionOpenAIAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("vision_openai", "Analyze Image (Vision)")

	getVisionModels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		models := []map[string]interface{}{
			{
				"id":   "gpt-4-vision-preview",
				"name": "GPT-4 Vision Preview (Latest)",
			},
			{
				"id":   "gpt-4-turbo",
				"name": "GPT-4 Turbo (With Vision)",
			},
			{
				"id":   "gpt-4o",
				"name": "GPT-4o (Optimized, With Vision)",
			},
			{
				"id":   "gpt-4o-mini",
				"name": "GPT-4o Mini (Faster, With Vision)",
			},
		}

		return ctx.Respond(models, len(models))
	}

	form.SelectField("model", "Model").
		Required(true).
		HelpText("Choose the vision model. GPT-4o is recommended for best performance.").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getVisionModels)).
				End().
				GetDynamicSource(),
		)

	form.SelectField("image_type", "Image Input Type").
		Required(true).
		HelpText("Choose how you'll provide the image").
		AddOptions([]*smartform.Option{
			{Value: "url", Label: "Image URL"},
			{Value: "base64", Label: "Base64 Encoded Image"},
		}...)

	form.TextareaField("image_input", "Image URL or Base64").
		Required(true).
		HelpText("Provide either a public image URL or base64 encoded image data (without data:image prefix)")

	form.TextareaField("prompt", "Analysis Prompt").
		Required(true).
		HelpText("What would you like to know about this image? Be specific for best results.")

	form.TextareaField("system_prompt", "System Prompt").
		Required(false).
		HelpText("Optional system message to guide the analysis behavior (e.g., 'You are an expert at reading invoices and extracting data')")

	form.SelectField("detail", "Analysis Detail Level").
		Required(false).
		HelpText("Control the detail level of image analysis. 'Low' is faster/cheaper, 'High' is more detailed.").
		AddOptions([]*smartform.Option{
			{Value: "auto", Label: "Auto (Default)"},
			{Value: "low", Label: "Low (Fast, 512x512)"},
			{Value: "high", Label: "High (Detailed, Full Resolution)"},
		}...)

	form.TextareaField("multiple_images", "Additional Images").
		Required(false).
		HelpText("Add more image URLs or base64 data for multi-image analysis. Separate with ||IMAGE||")

	form.NumberField("max_tokens", "Max Tokens").
		Required(false).
		HelpText("Maximum tokens for the response (default: 1000)")

	form.NumberField("temperature", "Temperature").
		Required(false).
		HelpText("Control response creativity (0 = focused, 2 = creative). Default: 0.3 for accuracy")

	schema := form.Build()
	return schema
}

func (a *VisionOpenAIAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *VisionOpenAIAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Parse input
	rawInput := ctx.Input()

	input := &visionOpenAIActionProps{
		Model:      rawInput["model"].(string),
		ImageInput: rawInput["image_input"].(string),
		ImageType:  rawInput["image_type"].(string),
		Prompt:     rawInput["prompt"].(string),
	}

	// Handle optional fields
	if systemPrompt, ok := rawInput["system_prompt"].(string); ok && systemPrompt != "" {
		input.SystemPrompt = systemPrompt
	}

	if detail, ok := rawInput["detail"].(string); ok && detail != "" {
		input.Detail = detail
	}

	if temp, ok := rawInput["temperature"].(float64); ok {
		input.Temperature = &temp
	}

	if maxTokens, ok := rawInput["max_tokens"].(float64); ok {
		maxTokensInt := int(maxTokens)
		input.MaxTokens = &maxTokensInt
	}

	// Parse multiple images if provided
	if multiImages, ok := rawInput["multiple_images"].(string); ok && multiImages != "" {
		input.MultipleImages = strings.Split(multiImages, "||IMAGE||")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Validate input
	if err := validateVisionInput(input); err != nil {
		return nil, err
	}

	// Build request body
	requestBody := buildVisionRequestBody(input)

	// Make API call
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
		return nil, ParseOpenAIError(res.Status().Code(), bodyBytes)
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
		return nil, errors.New("no analysis result returned")
	}

	// Build response
	response := map[string]interface{}{
		"vision_analysis": chatCompletion.Choices[0].Message.Content,
		"model":           input.Model,
		"prompt":          input.Prompt,
		"image_type":      input.ImageType,
		"detail_level":    input.Detail,
		"tokens_used":     chatCompletion.Usage.TotalTokens,
	}

	if input.SystemPrompt != "" {
		response["system_prompt"] = input.SystemPrompt
	}

	if len(input.MultipleImages) > 0 {
		response["images_analyzed"] = len(input.MultipleImages) + 1
	}

	return response, nil
}

func NewVisionOpenAIAction() sdk.Action {
	return &VisionOpenAIAction{}
}

// ============== Helper Functions ==============

func validateVisionInput(input *visionOpenAIActionProps) error {
	if input.Model == "" {
		return errors.New("model is required")
	}

	if input.ImageInput == "" {
		return errors.New("image input is required")
	}

	if input.ImageType != "url" && input.ImageType != "base64" {
		return errors.New("image type must be either 'url' or 'base64'")
	}

	if input.Prompt == "" {
		return errors.New("analysis prompt is required")
	}

	// Validate URL format if URL type
	if input.ImageType == "url" {
		if !strings.HasPrefix(input.ImageInput, "http://") && !strings.HasPrefix(input.ImageInput, "https://") {
			return errors.New("image URL must start with http:// or https://")
		}
	}

	// Validate detail level
	if input.Detail != "" && input.Detail != "low" && input.Detail != "high" && input.Detail != "auto" {
		return errors.New("detail must be 'low', 'high', or 'auto'")
	}

	if input.Temperature != nil && (*input.Temperature < 0 || *input.Temperature > 2) {
		return errors.New("temperature must be between 0 and 2")
	}

	return nil
}

func buildVisionRequestBody(input *visionOpenAIActionProps) map[string]interface{} {
	// Build messages array
	messages := []interface{}{}

	// Add system message if provided
	if input.SystemPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": input.SystemPrompt,
		})
	}

	// Build content array for user message
	content := []interface{}{
		map[string]interface{}{
			"type": "text",
			"text": input.Prompt,
		},
	}

	// Add main image
	imageContent := buildImageContent(input.ImageInput, input.ImageType, input.Detail)
	content = append(content, imageContent)

	// Add additional images if provided
	for _, additionalImage := range input.MultipleImages {
		additionalImage = strings.TrimSpace(additionalImage)
		if additionalImage != "" {
			imageContent := buildImageContent(additionalImage, input.ImageType, input.Detail)
			content = append(content, imageContent)
		}
	}

	// Add user message with all content
	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": content,
	})

	requestBody := map[string]interface{}{
		"model":    input.Model,
		"messages": messages,
	}

	// Add optional parameters
	if input.MaxTokens != nil {
		requestBody["max_tokens"] = *input.MaxTokens
	} else {
		requestBody["max_tokens"] = 1000 // Default for vision
	}

	if input.Temperature != nil {
		requestBody["temperature"] = *input.Temperature
	} else {
		requestBody["temperature"] = 0.3 // Lower default for accuracy
	}

	return requestBody
}

func buildImageContent(imageData string, imageType string, detail string) map[string]interface{} {
	imageContent := map[string]interface{}{
		"type": "image_url",
	}

	imageUrl := map[string]interface{}{}

	if imageType == "base64" {
		// Ensure proper base64 format
		if !strings.Contains(imageData, "data:image") {
			// Detect image format if possible, default to jpeg
			imageFormat := "jpeg"
			if strings.HasPrefix(imageData, "iVBORw0KGgo") {
				imageFormat = "png"
			} else if strings.HasPrefix(imageData, "R0lGOD") {
				imageFormat = "gif"
			}
			imageData = "data:image/" + imageFormat + ";base64," + imageData
		}
		imageUrl["url"] = imageData
	} else {
		imageUrl["url"] = imageData
	}

	// Add detail level if specified
	if detail != "" && detail != "auto" {
		imageUrl["detail"] = detail
	}

	imageContent["image_url"] = imageUrl

	return imageContent
}
