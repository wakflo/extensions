package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type analyzeImageActionProps struct {
	ImageURL string `json:"image_url"`
	Prompt   string `json:"prompt"`
	Model    string `json:"model"`
}

type AnalyzeImageAction struct{}

func (a *AnalyzeImageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "analyze_image_gemini",
		DisplayName:   "Analyze Image with Vision",
		Description:   "Analyze images using Gemini's vision capabilities to extract information, describe content, or answer questions about images.",
		Type:          core.ActionTypeAction,
		Documentation: analyzeImageGeminiDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"analysis": "The image shows a modern office space with...",
			"model":    "gemini-1.5-flash",
			"prompt":   "Describe what's in this image",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *AnalyzeImageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("analyze_image_gemini", "Analyze Image with Vision")

	form.TextField("image_url", "Image URL").
		Placeholder("https://example.com/image.jpg").
		HelpText("URL of the image to analyze (must be publicly accessible)").
		Required(true)

	form.TextareaField("prompt", "Analysis Prompt").
		Placeholder("Describe what's in this image").
		HelpText("What do you want to know about the image?").
		Required(true)

	RegisterVisionModelProps(form)

	return form.Build()
}

func (a *AnalyzeImageAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *AnalyzeImageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[analyzeImageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Validate image URL
	if input.ImageURL == "" {
		return nil, fmt.Errorf("image URL is required")
	}

	authCtx, err := ctx.AuthContext()
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
	model := client.GenerativeModel(modelName)

	// Create image part from URL
	imagePart := genai.FileData{
		URI: input.ImageURL,
	}

	// Generate content with image and prompt
	content, err := model.GenerateContent(gctx, imagePart, genai.Text(input.Prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	// Extract response text
	var response string
	if content != nil && len(content.Candidates) > 0 {
		for _, part := range content.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				response += string(text)
			}
		}
	}

	return map[string]interface{}{
		"analysis":  response,
		"model":     modelName,
		"prompt":    input.Prompt,
		"image_url": input.ImageURL,
	}, nil
}

// RegisterVisionModelProps registers only vision-capable models
func RegisterVisionModelProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	return form.SelectField("model", "Vision Model").
		Placeholder("Select a vision-capable model").
		Required(true).
		AddOptions(
			smartform.NewOption("gemini-1.5-flash", "Gemini 1.5 Flash"),
			smartform.NewOption("gemini-1.5-flash-latest", "Gemini 1.5 Flash Latest"),
			smartform.NewOption("gemini-1.5-pro", "Gemini 1.5 Pro"),
			smartform.NewOption("gemini-1.5-pro-latest", "Gemini 1.5 Pro Latest"),
			smartform.NewOption("gemini-2.0-flash", "Gemini 2.0 Flash"),
			smartform.NewOption("gemini-2.0-flash-exp", "Gemini 2.0 Flash Experimental"),
		).
		HelpText("Select a Gemini model with vision capabilities")
}

func NewAnalyzeImageAction() sdk.Action {
	return &AnalyzeImageAction{}
}
