package actions

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func CreateGeminiClient(ctx context.Context, auth *sdkcore.AuthContext) (*genai.Client, error) {
	return genai.NewClient(ctx, option.WithAPIKey(auth.Secret))
}

func GetModelInput() *sdkcore.AutoFormSchema {
	getModels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		gctx := context.Background()

		client, err := CreateGeminiClient(gctx, ctx.Auth)
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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Model").
		SetDescription("select gemini model").
		SetDynamicOptions(&getModels).
		SetRequired(true).Build()
}
