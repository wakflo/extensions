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

type embeddingsOpenAIActionProps struct {
	Model          string   `json:"model"`
	Input          []string `json:"input"`
	Dimensions     *int     `json:"dimensions,omitempty"`      // Only for text-embedding-3-* models
	EncodingFormat string   `json:"encoding_format,omitempty"` // "float" or "base64"
	User           string   `json:"user,omitempty"`            // Unique identifier for end-user
	ChunkSize      *int     `json:"chunk_size,omitempty"`      // For automatic chunking
	OverlapSize    *int     `json:"overlap_size,omitempty"`    // Overlap between chunks
}

type EmbeddingsResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// ============== Main Action ==============

type EmbeddingsOpenAIAction struct{}

// Metadata returns metadata about the action
func (a *EmbeddingsOpenAIAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "embeddings_openai",
		DisplayName:   "Create Embeddings",
		Description:   "Generate vector embeddings for text using OpenAI's embedding models. Perfect for semantic search, similarity matching, clustering, and building RAG (Retrieval-Augmented Generation) systems. Supports batch processing and automatic text chunking.",
		Type:          core.ActionTypeAction,
		Documentation: embeddingsOpenAIDocs,
		SampleOutput: map[string]any{
			"embeddings": []map[string]any{
				{
					"index":     0,
					"text":      "Sample text that was embedded",
					"embedding": []float64{0.0023, -0.0091, 0.0153},
					"tokens":    15,
				},
			},
			"model":          "text-embedding-3-small",
			"total_tokens":   15,
			"dimensions":     1536,
			"chunks_created": 1,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *EmbeddingsOpenAIAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("embeddings_openai", "Create Embeddings")

	getEmbeddingModels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		models := []map[string]interface{}{
			{
				"id":   "text-embedding-3-small",
				"name": "text-embedding-3-small (Latest, Cost-effective)",
			},
			{
				"id":   "text-embedding-3-large",
				"name": "text-embedding-3-large (Latest, High accuracy)",
			},
			{
				"id":   "text-embedding-ada-002",
				"name": "text-embedding-ada-002 (Legacy)",
			},
		}

		return ctx.Respond(models, len(models))
	}

	form.SelectField("model", "Embedding Model").
		Required(true).
		HelpText("Choose the embedding model. Use text-embedding-3-small for most cases, text-embedding-3-large for higher accuracy.").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getEmbeddingModels)).
				End().
				GetDynamicSource(),
		)

	form.TextareaField("input", "Text Input").
		Required(true).
		HelpText("Enter the text(s) to create embeddings for. For multiple texts, separate each with ||SEPARATOR||")

	form.NumberField("dimensions", "Vector Dimensions").
		Required(false).
		HelpText("Number of dimensions for the embedding vector (only for text-embedding-3-* models). Lower dimensions = faster/cheaper, higher = more accurate. Leave empty for model default.")

	form.NumberField("chunk_size", "Chunk Size (tokens)").
		Required(false).
		HelpText("If set, automatically split large texts into chunks of this token size (recommended: 500-2000 tokens)")

	form.NumberField("overlap_size", "Chunk Overlap (tokens)").
		Required(false).
		HelpText("Number of overlapping tokens between chunks to maintain context (recommended: 50-200 tokens)")

	form.SelectField("encoding_format", "Encoding Format").
		Required(false).
		HelpText("Format of the returned embeddings").
		AddOptions([]*smartform.Option{
			{Value: "float", Label: "Float Array (Default)"},
			{Value: "base64", Label: "Base64 Encoded"},
		}...)

	form.TextField("user", "User Identifier").
		Required(false).
		HelpText("Optional unique identifier for the end-user (for abuse monitoring)")

	schema := form.Build()
	return schema
}

func (a *EmbeddingsOpenAIAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *EmbeddingsOpenAIAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	rawInput := ctx.Input()
	inputText, _ := rawInput["input"].(string)

	var texts []string
	if strings.Contains(inputText, "||SEPARATOR||") {
		texts = strings.Split(inputText, "||SEPARATOR||")
	} else {
		texts = []string{inputText}
	}

	input := &embeddingsOpenAIActionProps{
		Model: rawInput["model"].(string),
		Input: texts,
	}

	if dim, ok := rawInput["dimensions"].(float64); ok {
		dimInt := int(dim)
		input.Dimensions = &dimInt
	}

	if chunkSize, ok := rawInput["chunk_size"].(float64); ok {
		chunkInt := int(chunkSize)
		input.ChunkSize = &chunkInt
	}

	if overlapSize, ok := rawInput["overlap_size"].(float64); ok {
		overlapInt := int(overlapSize)
		input.OverlapSize = &overlapInt
	}

	if format, ok := rawInput["encoding_format"].(string); ok && format != "" {
		input.EncodingFormat = format
	}

	if user, ok := rawInput["user"].(string); ok && user != "" {
		input.User = user
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if err := validateEmbeddingsInput(input); err != nil {
		return nil, err
	}

	processedTexts := input.Input
	if input.ChunkSize != nil && *input.ChunkSize > 0 {
		processedTexts = chunkTexts(input.Input, *input.ChunkSize, input.OverlapSize)
	}

	requestBody := buildEmbeddingsRequestBody(input, processedTexts)

	client, err := getOpenAiClient(authCtx.Extra["token"])
	if err != nil {
		return nil, err
	}

	res, err := client.POST("/embeddings").
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

	var embeddingsResp EmbeddingsResponse
	if err = json.Unmarshal(bodyBytes, &embeddingsResp); err != nil {
		return nil, errors.New("could not parse embeddings response")
	}

	// Format response
	embeddings := make([]map[string]interface{}, len(embeddingsResp.Data))
	for i, data := range embeddingsResp.Data {
		originalText := ""
		if i < len(processedTexts) {
			originalText = processedTexts[i]
			if len(originalText) > 100 {
				originalText = originalText[:100] + "..."
			}
		}

		embeddings[i] = map[string]interface{}{
			"index":     data.Index,
			"text":      originalText,
			"embedding": data.Embedding,
		}
	}

	dimensions := 1536 // default
	if input.Model == "text-embedding-3-small" && input.Dimensions == nil {
		dimensions = 1536
	} else if input.Model == "text-embedding-3-large" && input.Dimensions == nil {
		dimensions = 3072
	} else if input.Dimensions != nil {
		dimensions = *input.Dimensions
	}

	return map[string]interface{}{
		"embeddings":     embeddings,
		"model":          embeddingsResp.Model,
		"total_tokens":   embeddingsResp.Usage.TotalTokens,
		"dimensions":     dimensions,
		"chunks_created": len(processedTexts),
		"original_texts": len(input.Input),
	}, nil
}

func NewEmbeddingsOpenAIAction() sdk.Action {
	return &EmbeddingsOpenAIAction{}
}

func validateEmbeddingsInput(input *embeddingsOpenAIActionProps) error {
	if input.Model == "" {
		return errors.New("model is required")
	}

	if len(input.Input) == 0 {
		return errors.New("at least one input text is required")
	}

	if input.Dimensions != nil {
		if !strings.HasPrefix(input.Model, "text-embedding-3-") {
			return errors.New("dimensions can only be set for text-embedding-3-* models")
		}
		if *input.Dimensions < 1 || *input.Dimensions > 3072 {
			return errors.New("dimensions must be between 1 and 3072")
		}
	}

	if input.ChunkSize != nil && *input.ChunkSize < 100 {
		return errors.New("chunk size must be at least 100 tokens")
	}

	if input.OverlapSize != nil && input.ChunkSize != nil {
		if *input.OverlapSize >= *input.ChunkSize {
			return errors.New("overlap size must be smaller than chunk size")
		}
	}

	for i, text := range input.Input {
		if strings.TrimSpace(text) == "" {
			return fmt.Errorf("input text at index %d is empty", i)
		}
	}

	return nil
}

func buildEmbeddingsRequestBody(input *embeddingsOpenAIActionProps, processedTexts []string) map[string]interface{} {
	requestBody := map[string]interface{}{
		"model": input.Model,
		"input": processedTexts,
	}

	if input.Dimensions != nil {
		requestBody["dimensions"] = *input.Dimensions
	}

	if input.EncodingFormat != "" {
		requestBody["encoding_format"] = input.EncodingFormat
	}

	if input.User != "" {
		requestBody["user"] = input.User
	}

	return requestBody
}

func chunkTexts(texts []string, chunkSize int, overlapSize *int) []string {
	overlap := 0
	if overlapSize != nil {
		overlap = *overlapSize
	}

	var chunks []string
	for _, text := range texts {

		approxChunkChars := chunkSize * 4
		approxOverlapChars := overlap * 4

		if len(text) <= approxChunkChars {
			chunks = append(chunks, text)
			continue
		}

		for i := 0; i < len(text); i += (approxChunkChars - approxOverlapChars) {
			end := i + approxChunkChars
			if end > len(text) {
				end = len(text)
			}
			chunks = append(chunks, text[i:end])

			if end >= len(text) {
				break
			}
		}
	}

	return chunks
}
