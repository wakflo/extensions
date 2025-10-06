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

type summarizeTextActionProps struct {
	Text      string `json:"text"`
	Model     string `json:"model"`
	Style     string `json:"style"`
	MaxLength string `json:"max_length"`
}

type SummarizeTextAction struct{}

func (a *SummarizeTextAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:          "summarize_text_claude",
		DisplayName: "Summarize Text",
		Description: "Generate concise summaries of documents, articles, or any long-form content.",
		Type:        core.ActionTypeAction,
		Documentation: summarizeTetDocs,
		Icon: "",
		SampleOutput: map[string]any{
			"summary": "This document discusses...",
			"model":   "claude-3-5-sonnet-20241022",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SummarizeTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("summarize_text_claude", "Summarize Text")

	shared.RegisterModelProps(form)

	form.TextareaField("text", "Text to Summarize").
		Placeholder("Paste the text you want to summarize").
		HelpText("The content to summarize").
		Required(true)

	form.SelectField("style", "Summary Style").
		Placeholder("Select style").
		Required(false).
		AddOptions(
			smartform.NewOption("executive", "Executive Summary"),
			smartform.NewOption("bullets", "Bullet Points"),
			smartform.NewOption("abstract", "Academic Abstract"),
			smartform.NewOption("simple", "Simple Language"),
			smartform.NewOption("technical", "Technical Summary"),
		).
		HelpText("Style of summary")

	form.SelectField("max_length", "Length").
		Placeholder("Select length").
		Required(false).
		AddOptions(
			smartform.NewOption("brief", "Brief (1-2 paragraphs)"),
			smartform.NewOption("moderate", "Moderate (3-4 paragraphs)"),
			smartform.NewOption("detailed", "Detailed (5+ paragraphs)"),
		).
		HelpText("Desired summary length")

	return form.Build()
}

func (a *SummarizeTextAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *SummarizeTextAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[summarizeTextActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["apiKey"] == ""{
		return nil, errors.New("please add your claude api key to continue")
	}

	prompt := "Summarize the following text"
	
	switch input.Style {
	case "executive":
		prompt += " as an executive summary focusing on key business insights"
	case "bullets":
		prompt += " as clear bullet points"
	case "abstract":
		prompt += " as an academic abstract"
	case "simple":
		prompt += " in simple, easy-to-understand language"
	case "technical":
		prompt += " preserving technical details and terminology"
	}

	switch input.MaxLength {
	case "brief":
		prompt += " in 1-2 paragraphs"
	case "moderate":
		prompt += " in 3-4 paragraphs"
	case "detailed":
		prompt += " in 5 or more paragraphs with comprehensive detail"
	}

	prompt += ":\n\n" + input.Text

	request := shared.ClaudeRequest{
		Model: input.Model,
		Messages: []shared.ClaudeMessage{
			{
				Role: "user",
				Content: []interface{}{
					shared.ClaudeTextContent{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
		MaxTokens: 2048,
	}

	gctx := context.Background()
	response, err := shared.CallClaudeAPI(gctx, authCtx, request)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"summary": shared.ExtractResponseText(response),
		"model":   response.Model,
		"style":   input.Style,
	}, nil
}

func NewSummarizeTextAction() sdk.Action {
	return &SummarizeTextAction{}
}