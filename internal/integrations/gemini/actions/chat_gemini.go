package actions

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type chatGeminiActionProps struct {
	Chat  string `json:"chat" genai:"chat"`
	Model string `json:"model"`
}

type ChatGeminiAction struct{}

func (a *ChatGeminiAction) Name() string {
	return "Chat Gemini"
}

func (a *ChatGeminiAction) Description() string {
	return "ChatGPT: Seamlessly integrates with your workflow to enable AI-powered chatbots that can understand and respond to user queries, automating routine conversations and freeing up human agents to focus on complex tasks."
}

func (a *ChatGeminiAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ChatGeminiAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &chatGeminiDocs,
	}
}

func (a *ChatGeminiAction) Icon() *string {
	return nil
}

func (a *ChatGeminiAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"chat": autoform.NewLongTextField().
			SetLabel("Chat Prompt").
			SetRequired(true).
			SetPlaceholder("Enter your prompt here.").
			Build(),
		"model": GetModelInput(),
	}
}

func (a *ChatGeminiAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[chatGeminiActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	gctx := context.Background()
	client, err := CreateGeminiClient(gctx, ctx.Auth)
	if err != nil {
		return nil, err
	}

	sdkcore.PrettyPrint(input)

	content, err := client.GenerativeModel(input.Model).GenerateContent(gctx, genai.Text(input.Chat), nil)
	if err != nil {
		return nil, err
	}

	return content, err
}

func (a *ChatGeminiAction) Auth() *sdk.Auth {
	return nil
}

func (a *ChatGeminiAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ChatGeminiAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewChatGeminiAction() sdk.Action {
	return &ChatGeminiAction{}
}
