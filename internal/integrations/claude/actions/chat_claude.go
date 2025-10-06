package actions

import (
	"context"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/claude/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type chatClaudeActionProps struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	System      string  `json:"system"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type ChatClaudeAction struct{}

func (a *ChatClaudeAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "chat_claude",
		DisplayName:   "Chat with Claude",
		Description:   "Have a conversation with Claude AI for questions, analysis, creative writing, and problem-solving.",
		Type:          core.ActionTypeAction,
		Documentation: chatClaudeDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"response": "Hello! I'd be happy to help you...",
			"model":    "claude-3-5-sonnet-20241022",
			"usage": map[string]int{
				"input_tokens":  10,
				"output_tokens": 50,
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ChatClaudeAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("chat_claude", "Chat with Claude")

	shared.RegisterModelProps(form)

	form.TextareaField("prompt", "Message").
		Placeholder("Enter your message to Claude").
		HelpText("Your question or prompt").
		Required(true)

	form.TextareaField("system", "System Prompt").
		Placeholder("You are a helpful assistant...").
		HelpText("Optional system prompt to set Claude's behavior").
		Required(false)

	form.NumberField("temperature", "Temperature").
		Placeholder("0.7").
		HelpText("Controls randomness (0=focused, 1=creative)").
		Required(false)

	form.NumberField("max_tokens", "Max Tokens").
		Placeholder("1024").
		HelpText("Maximum response length in tokens").
		Required(false)

	return form.Build()
}

func (a *ChatClaudeAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *ChatClaudeAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[chatClaudeActionProps](ctx)
	if err != nil {
		return nil, err
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
	if input.Temperature == 0 {
		input.Temperature = 0.7
	}

	request := shared.ClaudeRequest{
		Model: input.Model,
		Messages: []shared.ClaudeMessage{
			{
				Role: "user",
				Content: []interface{}{
					shared.ClaudeTextContent{
						Type: "text",
						Text: input.Prompt,
					},
				},
			},
		},
		MaxTokens:   input.MaxTokens,
		Temperature: input.Temperature,
		System:      input.System,
	}

	gctx := context.Background()
	response, err := shared.CallClaudeAPI(gctx, authCtx, request)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"response": shared.ExtractResponseText(response),
		"model":    response.Model,
		"usage": map[string]int{
			"input_tokens":  response.Usage.InputTokens,
			"output_tokens": response.Usage.OutputTokens,
			"total_tokens":  response.Usage.InputTokens + response.Usage.OutputTokens,
		},
		"stop_reason": response.StopReason,
	}, nil
}

func NewChatClaudeAction() sdk.Action {
	return &ChatClaudeAction{}
}
