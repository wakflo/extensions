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

type generateEmbeddingActionProps struct {
	Text  string `json:"text"`
	Model string `json:"model"`
	Task  string `json:"task"`
}

type GenerateEmbeddingAction struct{}

// Metadata returns metadata about the action
func (a *GenerateEmbeddingAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "generate_embedding_gemini",
		DisplayName:   "Generate Text Embedding",
		Description:   "Generate embeddings for text using Gemini's embedding models. Useful for semantic search, clustering, and similarity comparisons.",
		Type:          core.ActionTypeAction,
		Documentation: embeddingsGeminiDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"embedding": []float32{0.1, 0.2, 0.3},
			"dimension": 768,
			"model":     "text-embedding-004",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GenerateEmbeddingAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("generate_embedding_gemini", "Generate Text Embedding")

	form.TextareaField("text", "Text").
		Placeholder("Enter text to generate embeddings for").
		HelpText("The text content to embed").
		Required(true)

	RegisterEmbeddingModelProps(form)

	form.SelectField("task", "Task Type").
		Placeholder("Select embedding task type").
		Required(false).
		AddOptions(
			smartform.NewOption("retrieval_document", "Retrieval Document"),
			smartform.NewOption("retrieval_query", "Retrieval Query"),
			smartform.NewOption("semantic_similarity", "Semantic Similarity"),
			smartform.NewOption("classification", "Classification"),
			smartform.NewOption("clustering", "Clustering"),
		).
		HelpText("The type of task for optimized embeddings (optional)")

	return form.Build()
}

func (a *GenerateEmbeddingAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *GenerateEmbeddingAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[generateEmbeddingActionProps](ctx)
	if err != nil {
		return nil, err
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

	modelName := strings.TrimPrefix(input.Model, "models/")

	embeddingModel := client.EmbeddingModel(modelName)

	// Set task type if provided
	if input.Task != "" {
		taskType := genai.TaskTypeUnspecified
		switch input.Task {
		case "retrieval_document":
			taskType = genai.TaskTypeRetrievalDocument
		case "retrieval_query":
			taskType = genai.TaskTypeRetrievalQuery
		case "semantic_similarity":
			taskType = genai.TaskTypeSemanticSimilarity
		case "classification":
			taskType = genai.TaskTypeClassification
		case "clustering":
			taskType = genai.TaskTypeClustering
		}
		embeddingModel.TaskType = taskType
	}

	resp, err := embeddingModel.EmbedContent(gctx, genai.Text(input.Text))
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	return map[string]interface{}{
		"embedding": resp.Embedding.Values,
		"dimension": len(resp.Embedding.Values),
		"model":     modelName,
		"task":      input.Task,
	}, nil
}

func NewGenerateEmbeddingAction() sdk.Action {
	return &GenerateEmbeddingAction{}
}
