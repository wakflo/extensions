package actions

import (
	"context"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/juicycleff/smartform/v1"
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
		ID:            "summarize_text_gemini",
		DisplayName:   "Summarize Text",
		Description:   "Generate concise summaries of long documents, articles, or conversations with customizable styles.",
		Type:          core.ActionTypeAction,
		Documentation: summarizeTextGeminiDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"summary": "This document discusses...",
			"model":   "gemini-1.5-flash",
			"style":   "executive",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SummarizeTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("summarize_text_gemini", "Summarize Text")

	form.TextareaField("text", "Text to Summarize").
		Placeholder("Paste the text you want to summarize").
		HelpText("The document or text to summarize").
		Required(true)

	form.SelectField("style", "Summary Style").
		Placeholder("Select summary style").
		Required(false).
		AddOptions(
			smartform.NewOption("executive", "Executive Summary"),
			smartform.NewOption("technical", "Technical Summary"),
			smartform.NewOption("bullets", "Bullet Points"),
			smartform.NewOption("abstract", "Abstract"),
			smartform.NewOption("simple", "Simple Language"),
		).
		HelpText("The style of summary to generate")

	form.SelectField("max_length", "Summary Length").
		Placeholder("Select maximum length").
		Required(false).
		AddOptions(
			smartform.NewOption("brief", "Brief (1-2 paragraphs)"),
			smartform.NewOption("moderate", "Moderate (3-4 paragraphs)"),
			smartform.NewOption("detailed", "Detailed (5+ paragraphs)"),
		).
		HelpText("Approximate length of the summary")

	RegisterModelProps(form)

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

	gctx := context.Background()
	client, err := CreateGeminiClient(gctx, authCtx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	modelName := strings.TrimPrefix(input.Model, "models/")
	model := client.GenerativeModel(modelName)

	// Build prompt based on style and length
	var prompt string
	switch input.Style {
	case "executive":
		prompt = "Provide an executive summary of the following text. Focus on key business implications and strategic insights:\n\n"
	case "technical":
		prompt = "Provide a detailed technical summary of the following text, preserving important technical details:\n\n"
	case "bullets":
		prompt = "Summarize the following text as bullet points. List the main ideas and key findings:\n\n"
	case "abstract":
		prompt = "Write an academic-style abstract for the following text:\n\n"
	case "simple":
		prompt = "Summarize the following text in simple, plain language that anyone can understand:\n\n"
	default:
		prompt = "Summarize the following text concisely:\n\n"
	}

	// Add length instruction
	switch input.MaxLength {
	case "brief":
		prompt += "(Keep the summary to 1-2 paragraphs)\n\n"
	case "moderate":
		prompt += "(Provide a moderate summary of 3-4 paragraphs)\n\n"
	case "detailed":
		prompt += "(Provide a detailed summary of 5 or more paragraphs)\n\n"
	}

	prompt += input.Text

	content, err := model.GenerateContent(gctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	var summary string
	if content != nil && len(content.Candidates) > 0 {
		for _, part := range content.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				summary += string(text)
			}
		}
	}

	return map[string]interface{}{
		"summary": summary,
		"model":   modelName,
		"style":   input.Style,
	}, nil
}

func NewSummarizeTextAction() sdk.Action {
	return &SummarizeTextAction{}
}
