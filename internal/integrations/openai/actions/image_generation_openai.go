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
)

// ============== Type Definitions ==============

type imageGenerationOpenAIActionProps struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt,omitempty"` // Things to avoid
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	Style          string `json:"style,omitempty"`
	NumImages      *int   `json:"n,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"` // "url" or "b64_json"
	User           string `json:"user,omitempty"`
}

type ImageGenerationResponses struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL           string `json:"url,omitempty"`
		B64JSON       string `json:"b64_json,omitempty"`
		RevisedPrompt string `json:"revised_prompt,omitempty"`
	} `json:"data"`
}

type ImageGenerationOpenAIAction struct{}

// Metadata returns metadata about the action
func (a *ImageGenerationOpenAIAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:          "image_generation_openai",
		DisplayName: "Generate Image (DALL-E)",
		Description: "Create images from text descriptions using OpenAI's DALL-E models. Generate product mockups, marketing materials, illustrations, concept art, or any visual content from natural language prompts.",
		Type:        core.ActionTypeAction,
		Documentation: `
# OpenAI Image Generation Action

Generate images from text descriptions using DALL-E 3 or DALL-E 2.

## Key Features
- **Text to Image**: Convert natural language descriptions into images
- **Multiple Styles**: Choose between vivid or natural styles (DALL-E 3)
- **High Quality**: Generate HD quality images for professional use
- **Various Sizes**: Support for different aspect ratios and resolutions
- **Batch Generation**: Create multiple variations at once (DALL-E 2)
- **Safety Filters**: Built-in content moderation

## Use Cases

### Marketing & Content
- Social media graphics
- Blog post illustrations
- Ad creatives
- Email headers
- Banner images

### E-commerce
- Product mockups
- Category images
- Promotional graphics
- Placeholder images

### Design & Creative
- Concept art
- Mood boards
- Storyboards
- Character designs
- Logo concepts

### Documentation
- Tutorial illustrations
- Process diagrams
- Infographic elements
- Technical illustrations

## Model Differences

### DALL-E 3 (Recommended)
- Better prompt understanding
- More accurate text rendering
- Higher quality output
- Supports vivid/natural styles
- One image per request
- Sizes: 1024x1024, 1024x1792, 1792x1024

### DALL-E 2
- Faster generation
- Lower cost
- Can generate up to 10 variations
- Sizes: 256x256, 512x512, 1024x1024

## Prompt Tips
- Be specific and detailed
- Include style references (e.g., "oil painting", "3D render")
- Describe composition and perspective
- Mention lighting and mood
- Specify colors and textures

## Best Practices
- Use DALL-E 3 for production quality
- Include negative prompts to avoid unwanted elements
- Test prompts iteratively
- Save successful prompts as templates
- Consider copyright and usage rights
`,
		SampleOutput: map[string]any{
			"images": []map[string]any{
				{
					"url":            "https://oaidalleapi.../img-ABC123.png",
					"revised_prompt": "A modern minimalist office space with a standing desk, dual monitors displaying analytics dashboards, indoor plants, large windows with city view, warm natural lighting, photorealistic style",
					"index":          0,
				},
			},
			"model":      "dall-e-3",
			"prompt":     "Modern minimalist office with analytics dashboards",
			"size":       "1792x1024",
			"quality":    "hd",
			"style":      "natural",
			"created_at": "2024-01-20T10:30:00Z",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ImageGenerationOpenAIAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("image_generation_openai", "Generate Image (DALL-E)")

	form.SelectField("model", "Model").
		Required(true).
		HelpText("Choose the DALL-E model. DALL-E 3 provides better quality and understanding.").
		AddOptions([]*smartform.Option{
			{Value: "dall-e-3", Label: "DALL-E 3 (Best Quality, Latest)"},
			{Value: "dall-e-2", Label: "DALL-E 2 (Faster, Lower Cost)"},
		}...)

	form.TextareaField("prompt", "Image Prompt").
		Required(true).
		HelpText("Describe the image you want to create. Be specific about style, composition, colors, and mood for best results.")

	form.TextareaField("negative_prompt", "Negative Prompt").
		Required(false).
		HelpText("Optional: Describe what you DON'T want in the image (e.g., 'no text, no watermarks, no people')")

	form.SelectField("size", "Image Size").
		Required(false).
		HelpText("Choose image dimensions. DALL-E 3 supports wide/tall formats, DALL-E 2 only square.").
		AddOptions([]*smartform.Option{
			{Value: "1024x1024", Label: "1024x1024 (Square, Default)"},
			{Value: "1792x1024", Label: "1792x1024 (Wide, DALL-E 3 only)"},
			{Value: "1024x1792", Label: "1024x1792 (Tall, DALL-E 3 only)"},
			{Value: "512x512", Label: "512x512 (Small, DALL-E 2 only)"},
			{Value: "256x256", Label: "256x256 (Tiny, DALL-E 2 only)"},
		}...)

	form.SelectField("quality", "Image Quality").
		Required(false).
		HelpText("Quality setting for DALL-E 3. HD takes longer but produces better results.").
		AddOptions([]*smartform.Option{
			{Value: "standard", Label: "Standard (Faster)"},
			{Value: "hd", Label: "HD (Best Quality, DALL-E 3 only)"},
		}...)

	form.SelectField("style", "Image Style").
		Required(false).
		HelpText("Visual style for DALL-E 3. Vivid gives more dramatic/hyper-real images.").
		AddOptions([]*smartform.Option{
			{Value: "natural", Label: "Natural (Realistic)"},
			{Value: "vivid", Label: "Vivid (Dramatic, Hyper-real)"},
		}...)

	form.NumberField("n", "Number of Images").
		Required(false).
		HelpText("Number of images to generate (1-10). DALL-E 3 only supports 1 image per request.")

	form.SelectField("response_format", "Response Format").
		Required(false).
		HelpText("How to receive the generated images.").
		AddOptions([]*smartform.Option{
			{Value: "url", Label: "URL (Default, expires after 1 hour)"},
			{Value: "b64_json", Label: "Base64 (For immediate storage)"},
		}...)

	form.TextField("user", "User Identifier").
		Required(false).
		HelpText("Optional unique identifier for the end-user (for abuse monitoring)")

	schema := form.Build()
	return schema
}

func (a *ImageGenerationOpenAIAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *ImageGenerationOpenAIAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	rawInput := ctx.Input()

	input := &imageGenerationOpenAIActionProps{
		Model:  rawInput["model"].(string),
		Prompt: rawInput["prompt"].(string),
	}

	// Handle optional fields
	if negPrompt, ok := rawInput["negative_prompt"].(string); ok && negPrompt != "" {
		input.Prompt = input.Prompt + ". Avoid: " + negPrompt
	}

	if size, ok := rawInput["size"].(string); ok && size != "" {
		input.Size = size
	}

	if quality, ok := rawInput["quality"].(string); ok && quality != "" {
		input.Quality = quality
	}

	if style, ok := rawInput["style"].(string); ok && style != "" {
		input.Style = style
	}

	if numImages, ok := rawInput["n"].(float64); ok {
		n := int(numImages)
		input.NumImages = &n
	}

	if format, ok := rawInput["response_format"].(string); ok && format != "" {
		input.ResponseFormat = format
	}

	if user, ok := rawInput["user"].(string); ok && user != "" {
		input.User = user
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Validate input
	if err := validateImageGenInput(input); err != nil {
		return nil, err
	}

	// Build request body
	requestBody := buildImageGenRequestBody(input)

	// Make API call
	client, err := getOpenAiClient(authCtx.Extra["token"])
	if err != nil {
		return nil, err
	}

	res, err := client.POST("/images/generations").
		Header().AddContentType("application/json").
		Body().AsJSON(requestBody).
		Send()
	if err != nil {
		return nil, err
	}

	if res.Status().IsError() {
		return nil, errors.New(res.Status().Text())
	}

	bodyBytes, err := io.ReadAll(res.Body().Raw())
	if err != nil {
		return nil, err
	}

	var imageResponse ImageGenerationResponses
	if err = json.Unmarshal(bodyBytes, &imageResponse); err != nil {
		return nil, errors.New("could not parse image generation response")
	}

	if len(imageResponse.Data) == 0 {
		return nil, errors.New("no images were generated")
	}

	// Simple response - just success message and URL as plain text
	imageURL := ""
	if imageResponse.Data[0].URL != "" {
		imageURL = imageResponse.Data[0].URL
	}
	fmt.Println(imageURL)

	safe := strings.Replace(imageURL, "https://", "https[:]//", 1)

	response := map[string]interface{}{
		"status":                 "success",
		"message":                "Image generated successfully! Copy the image_url below and paste in a new browser tab to view.",
		"image_url":              imageURL,
		"image_url_display_only": safe,
		"expires_in":             "1 hour",
		"prompt":                 rawInput["prompt"].(string),
		"model":                  input.Model,
	}

	if input.Size != "" {
		response["size"] = input.Size
	}
	if imageResponse.Data[0].RevisedPrompt != "" {
		response["revised_prompt"] = imageResponse.Data[0].RevisedPrompt
	}

	return response, nil
}

func NewImageGenerationOpenAIAction() sdk.Action {
	return &ImageGenerationOpenAIAction{}
}

func validateImageGenInput(input *imageGenerationOpenAIActionProps) error {
	if input.Model == "" {
		return errors.New("model is required")
	}

	if input.Prompt == "" {
		return errors.New("image prompt is required")
	}

	if input.Model == "dall-e-3" {
		// DALL-E 3 only supports n=1
		if input.NumImages != nil && *input.NumImages != 1 {
			return errors.New("DALL-E 3 only supports generating 1 image at a time")
		}

		// Validate size for DALL-E 3
		if input.Size != "" {
			validSizes := []string{"1024x1024", "1792x1024", "1024x1792"}
			valid := false
			for _, s := range validSizes {
				if input.Size == s {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid size for DALL-E 3: %s. Must be one of: 1024x1024, 1792x1024, 1024x1792", input.Size)
			}
		}

		if input.Quality != "" && input.Quality != "standard" && input.Quality != "hd" {
			return errors.New("quality must be 'standard' or 'hd'")
		}

		if input.Style != "" && input.Style != "natural" && input.Style != "vivid" {
			return errors.New("style must be 'natural' or 'vivid'")
		}
	}

	if input.Model == "dall-e-2" {
		if input.NumImages != nil && (*input.NumImages < 1 || *input.NumImages > 10) {
			return errors.New("number of images must be between 1 and 10")
		}

		if input.Size != "" {
			validSizes := []string{"256x256", "512x512", "1024x1024"}
			valid := false
			for _, s := range validSizes {
				if input.Size == s {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid size for DALL-E 2: %s. Must be one of: 256x256, 512x512, 1024x1024", input.Size)
			}
		}

		// These are DALL-E 3 only features
		if input.Quality != "" {
			return errors.New("quality setting is only available for DALL-E 3")
		}
		if input.Style != "" {
			return errors.New("style setting is only available for DALL-E 3")
		}
	}

	// Validate response format
	if input.ResponseFormat != "" && input.ResponseFormat != "url" && input.ResponseFormat != "b64_json" {
		return errors.New("response_format must be 'url' or 'b64_json'")
	}

	// Check prompt length (DALL-E has a 4000 character limit)
	if len(input.Prompt) > 4000 {
		return errors.New("prompt exceeds 4000 character limit")
	}

	return nil
}

func buildImageGenRequestBody(input *imageGenerationOpenAIActionProps) map[string]interface{} {
	requestBody := map[string]interface{}{
		"prompt": input.Prompt,
	}

	// Add model only for DALL-E 3
	if input.Model == "dall-e-3" {
		requestBody["model"] = input.Model
	}

	// Add size with defaults
	if input.Size != "" {
		requestBody["size"] = input.Size
	} else {
		requestBody["size"] = "1024x1024"
	}

	if input.NumImages != nil {
		requestBody["n"] = *input.NumImages
	} else {
		requestBody["n"] = 1
	}

	if input.ResponseFormat != "" {
		requestBody["response_format"] = input.ResponseFormat
	} else {
		requestBody["response_format"] = "url"
	}

	if input.Model == "dall-e-3" {
		if input.Quality != "" {
			requestBody["quality"] = input.Quality
		} else {
			requestBody["quality"] = "standard"
		}

		if input.Style != "" {
			requestBody["style"] = input.Style
		} else {
			requestBody["style"] = "natural"
		}
	}

	if input.User != "" {
		requestBody["user"] = input.User
	}

	return requestBody
}
