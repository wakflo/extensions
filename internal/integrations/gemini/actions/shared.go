package actions

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/juicycleff/smartform/v1"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

func CreateGeminiClient(ctx context.Context, auth *sdkcontext.AuthContext) (*genai.Client, error) {
	return genai.NewClient(ctx, option.WithAPIKey(auth.Extra["key"]))
}

func RegisterModelProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getModels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		gctx := context.Background()

		client, err := CreateGeminiClient(gctx, authCtx)
		if err != nil {
			return nil, err
		}

		iter := client.ListModels(gctx)

		var models []map[string]string
		for {
			model, err := iter.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
				return nil, err
			}
			models = append(models, map[string]string{
				"id":   model.Name,
				"name": model.DisplayName,
			})
		}

		fmt.Println("models>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", models)

		return ctx.Respond(models, len(models))
	}

	return form.SelectField("model", "Model").
		Placeholder("Select a Gemini model").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getModels)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a Gemini model")
}

func RegisterEmbeddingModelProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	return form.SelectField("model", "Embedding Model").
		Placeholder("Select an embedding model").
		Required(true).
		AddOptions(
			smartform.NewOption("models/text-embedding-004", "Text Embedding 004"),
			smartform.NewOption("models/embedding-001", "Embedding 001"),
			smartform.NewOption("models/text-embedding-preview-0815", "Text Embedding Preview"),
		).
		HelpText("Select a Gemini embedding model")
}
