package actions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/claude/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type chatWithImagesActionProps struct {
	Prompt    string `json:"prompt"`
	ImageURL  string `json:"image_url"`
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
}

type ChatWithImagesAction struct{}

func (a *ChatWithImagesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "chat_with_images_claude",
		DisplayName:   "Analyze Images with Claude",
		Description:   "Analyze images and answer questions about visual content using Claude's vision capabilities.",
		Type:          core.ActionTypeAction,
		Documentation: chatImageDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"response": "The image shows...",
			"model":    "claude-3-5-sonnet-20241022",
			"usage": map[string]int{
				"input_tokens":  100,
				"output_tokens": 200,
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ChatWithImagesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("chat_with_images_claude", "Analyze Images")

	shared.RegisterModelProps(form)

	form.TextareaField("prompt", "Question").
		Placeholder("What's in this image?").
		HelpText("What do you want to know about the image?").
		Required(true)

	form.TextField("image_url", "Image URL").
		Placeholder("https://example.com/image.jpg").
		HelpText("URL of the image to analyze (must be publicly accessible)").
		Required(true)

	form.NumberField("max_tokens", "Max Tokens").
		Placeholder("1024").
		HelpText("Maximum response length").
		Required(false)

	return form.Build()
}

func (a *ChatWithImagesAction) Auth() *core.AuthMetadata {
	return nil
}

// fetchImageFromURL downloads an image from a URL and returns its bytes
func fetchImageFromURL(url string) ([]byte, string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch image: HTTP %d", resp.StatusCode)
	}

	// Read image data
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Check if the response is too large (Claude has a limit of ~5MB for images)
	if len(imageBytes) > 5*1024*1024 {
		return nil, "", fmt.Errorf("image is too large (max 5MB)")
	}

	mimeType := "image/jpeg" // default
	if len(imageBytes) > 4 {
		if imageBytes[0] == 0x89 && imageBytes[1] == 0x50 && imageBytes[2] == 0x4E && imageBytes[3] == 0x47 {
			mimeType = "image/png"
		} else if imageBytes[0] == 0x47 && imageBytes[1] == 0x49 && imageBytes[2] == 0x46 {
			mimeType = "image/gif"
		} else if string(imageBytes[:4]) == "RIFF" && len(imageBytes) > 11 && string(imageBytes[8:12]) == "WEBP" {
			mimeType = "image/webp"
		} else if imageBytes[0] == 0xFF && imageBytes[1] == 0xD8 {
			mimeType = "image/jpeg"
		}
	}

	// Also try to get MIME type from Content-Type header as fallback
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		switch contentType {
		case "image/png":
			mimeType = "image/png"
		case "image/gif":
			mimeType = "image/gif"
		case "image/webp":
			mimeType = "image/webp"
		case "image/jpeg", "image/jpg":
			mimeType = "image/jpeg"
		}
	}

	return imageBytes, mimeType, nil
}

func (a *ChatWithImagesAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[chatWithImagesActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.ImageURL == "" {
		return nil, errors.New("image URL is required")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["apiKey"] == "" {
		return nil, errors.New("please add your claude api key to continue")
	}

	if input.Model == "" {
		return nil, errors.New("Model is required")
	}

	if input.MaxTokens == 0 {
		input.MaxTokens = 1024
	}

	imageBytes, mimeType, err := fetchImageFromURL(input.ImageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}

	base64Data := base64.StdEncoding.EncodeToString(imageBytes)

	content := []interface{}{
		shared.ClaudeImageContent{
			Type: "image",
			Source: shared.ClaudeImageSource{
				Type:      "base64",
				MediaType: mimeType,
				Data:      base64Data,
			},
		},
		shared.ClaudeTextContent{
			Type: "text",
			Text: input.Prompt,
		},
	}

	request := shared.ClaudeRequest{
		Model: input.Model,
		Messages: []shared.ClaudeMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
		MaxTokens: input.MaxTokens,
	}

	gctx := context.Background()
	response, err := shared.CallClaudeAPI(gctx, authCtx, request)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"response":  shared.ExtractResponseText(response),
		"model":     response.Model,
		"image_url": input.ImageURL,
		"usage": map[string]int{
			"input_tokens":  response.Usage.InputTokens,
			"output_tokens": response.Usage.OutputTokens,
			"total_tokens":  response.Usage.InputTokens + response.Usage.OutputTokens,
		},
	}, nil
}

func NewChatWithImagesAction() sdk.Action {
	return &ChatWithImagesAction{}
}
