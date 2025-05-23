package actions

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type chatGeminiActionProps struct {
	Chat  string `json:"chat" genai:"chat"`
	Model string `json:"model"`
}

type ChatGeminiAction struct{}

// Metadata returns metadata about the action
func (a *ChatGeminiAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "chat_gemini",
		DisplayName:   "Chat Gemini",
		Description:   "ChatGPT: Seamlessly integrates with your workflow to enable AI-powered chatbots that can understand and respond to user queries, automating routine conversations and freeing up human agents to focus on complex tasks.",
		Type:          core.ActionTypeAction,
		Documentation: chatGeminiDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ChatGeminiAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("chat_gemini", "Chat Gemini")

	form.TextareaField("chat", "chat").
		Placeholder("Enter your prompt here.").
		HelpText("Chat Prompt").
		Required(true)

	RegisterModelProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ChatGeminiAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ChatGeminiAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[chatGeminiActionProps](ctx)
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

	// core.PrettyPrint(input)

	content, err := client.GenerativeModel(input.Model).GenerateContent(gctx, genai.Text(input.Chat), nil)
	if err != nil {
		return nil, err
	}

	return content, err
}

func NewChatGeminiAction() sdk.Action {
	return &ChatGeminiAction{}
}
