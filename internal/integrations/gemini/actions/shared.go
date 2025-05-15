package actions

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/juicycleff/smartform/v1"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

func CreateGeminiClient(ctx context.Context, auth *sdkcore.AuthContext) (*genai.Client, error) {
	return genai.NewClient(ctx, option.WithAPIKey(auth.Secret))
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
				"id":   model.BaseModelID,
				"name": model.DisplayName,
			})
		}

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
				End().
				GetDynamicSource(),
		).
		HelpText("Select a Gemini model")
}
