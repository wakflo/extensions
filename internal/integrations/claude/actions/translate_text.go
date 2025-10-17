// translate_text.go - Translation action for Claude integration

package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/claude/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type translateTextActionProps struct {
	Text         string  `json:"text"`
	TargetLang   string  `json:"target_lang"`
	SourceLang   string  `json:"source_lang"`
	Model        string  `json:"model"`
	Formality    string  `json:"formality"`
	PreserveTone bool    `json:"preserve_tone"`
	Temperature  float64 `json:"temperature"`
}

type TranslateTextAction struct{}

func (a *TranslateTextAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "translate_text_claude",
		DisplayName:   "Translate Text",
		Description:   "Translate text between multiple languages with context awareness and cultural adaptation using Claude.",
		Type:          core.ActionTypeAction,
		Documentation: translateTextDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"translated_text":   "Bonjour le monde!",
			"source_language":   "English",
			"target_language":   "French",
			"detected_language": "en",
			"confidence":        0.95,
			"alternative":       "Salut le monde!",
			"model":             "claude-3-5-sonnet-20241022",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *TranslateTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("translate_text_claude", "Translate Text")

	shared.RegisterModelProps(form)

	form.TextareaField("text", "Text to Translate").
		Placeholder("Enter the text you want to translate...").
		HelpText("The text content to translate (up to 100,000 characters)").
		Required(true)

	form.SelectField("target_lang", "Target Language").
		Placeholder("Select target language").
		Required(true).
		AddOptions(
			// Major Languages
			smartform.NewOption("en", "English"),
			smartform.NewOption("es", "Spanish"),
			smartform.NewOption("fr", "French"),
			smartform.NewOption("de", "German"),
			smartform.NewOption("it", "Italian"),
			smartform.NewOption("pt", "Portuguese"),
			smartform.NewOption("pt-BR", "Portuguese (Brazilian)"),
			smartform.NewOption("ru", "Russian"),
			smartform.NewOption("ja", "Japanese"),
			smartform.NewOption("ko", "Korean"),
			smartform.NewOption("zh", "Chinese (Simplified)"),
			smartform.NewOption("zh-TW", "Chinese (Traditional)"),
			smartform.NewOption("ar", "Arabic"),
			smartform.NewOption("hi", "Hindi"),
			// European Languages
			smartform.NewOption("nl", "Dutch"),
			smartform.NewOption("pl", "Polish"),
			smartform.NewOption("sv", "Swedish"),
			smartform.NewOption("no", "Norwegian"),
			smartform.NewOption("da", "Danish"),
			smartform.NewOption("fi", "Finnish"),
			smartform.NewOption("el", "Greek"),
			smartform.NewOption("cs", "Czech"),
			smartform.NewOption("hu", "Hungarian"),
			smartform.NewOption("ro", "Romanian"),
			smartform.NewOption("bg", "Bulgarian"),
			smartform.NewOption("uk", "Ukrainian"),
			smartform.NewOption("hr", "Croatian"),
			smartform.NewOption("sr", "Serbian"),
			smartform.NewOption("sk", "Slovak"),
			smartform.NewOption("sl", "Slovenian"),
			// Asian Languages
			smartform.NewOption("th", "Thai"),
			smartform.NewOption("vi", "Vietnamese"),
			smartform.NewOption("id", "Indonesian"),
			smartform.NewOption("ms", "Malay"),
			smartform.NewOption("tr", "Turkish"),
			smartform.NewOption("he", "Hebrew"),
			smartform.NewOption("fa", "Persian (Farsi)"),
			smartform.NewOption("ur", "Urdu"),
			smartform.NewOption("bn", "Bengali"),
			smartform.NewOption("ta", "Tamil"),
			smartform.NewOption("te", "Telugu"),
			smartform.NewOption("mr", "Marathi"),
			smartform.NewOption("gu", "Gujarati"),
			smartform.NewOption("kn", "Kannada"),
			smartform.NewOption("ml", "Malayalam"),
			// Other Languages
			smartform.NewOption("fil", "Filipino"),
			smartform.NewOption("sw", "Swahili"),
			smartform.NewOption("af", "Afrikaans"),
			smartform.NewOption("ca", "Catalan"),
			smartform.NewOption("eu", "Basque"),
			smartform.NewOption("ga", "Irish"),
			smartform.NewOption("cy", "Welsh"),
			smartform.NewOption("is", "Icelandic"),
			smartform.NewOption("et", "Estonian"),
			smartform.NewOption("lv", "Latvian"),
			smartform.NewOption("lt", "Lithuanian"),
			smartform.NewOption("mt", "Maltese"),
			smartform.NewOption("sq", "Albanian"),
			smartform.NewOption("mk", "Macedonian"),
		).
		HelpText("The language to translate the text into")

	form.SelectField("source_lang", "Source Language").
		Placeholder("Auto-detect or select").
		Required(false).
		AddOptions(
			smartform.NewOption("auto", "Auto-detect"),
			smartform.NewOption("en", "English"),
			smartform.NewOption("es", "Spanish"),
			smartform.NewOption("fr", "French"),
			smartform.NewOption("de", "German"),
			smartform.NewOption("it", "Italian"),
			smartform.NewOption("pt", "Portuguese"),
			smartform.NewOption("ru", "Russian"),
			smartform.NewOption("ja", "Japanese"),
			smartform.NewOption("ko", "Korean"),
			smartform.NewOption("zh", "Chinese (Simplified)"),
			smartform.NewOption("zh-TW", "Chinese (Traditional)"),
			smartform.NewOption("ar", "Arabic"),
			smartform.NewOption("hi", "Hindi"),
			smartform.NewOption("nl", "Dutch"),
			smartform.NewOption("pl", "Polish"),
			smartform.NewOption("tr", "Turkish"),
			smartform.NewOption("vi", "Vietnamese"),
			smartform.NewOption("th", "Thai"),
			smartform.NewOption("id", "Indonesian"),
			smartform.NewOption("he", "Hebrew"),
		).
		HelpText("The source language (leave empty for auto-detection)")

	form.SelectField("formality", "Formality Level").
		Placeholder("Select formality").
		Required(false).
		AddOptions(
			smartform.NewOption("auto", "Automatic"),
			smartform.NewOption("formal", "Formal"),
			smartform.NewOption("informal", "Informal"),
			smartform.NewOption("business", "Business"),
			smartform.NewOption("casual", "Casual"),
			smartform.NewOption("academic", "Academic"),
		).
		HelpText("The formality level of the translation")

	form.CheckboxField("preserve_tone", "Preserve Tone").
		HelpText("Maintain the original tone and style of the text").
		Required(false)

	form.NumberField("temperature", "Temperature").
		Placeholder("0.3").
		HelpText("Controls translation creativity (0=conservative, 1=creative, default: 0.3)").
		Required(false)

	return form.Build()
}

func (a *TranslateTextAction) Auth() *core.AuthMetadata {
	return nil
}

func getLanguageName(code string) string {
	languages := map[string]string{
		"en":    "English",
		"es":    "Spanish",
		"fr":    "French",
		"de":    "German",
		"it":    "Italian",
		"pt":    "Portuguese",
		"pt-BR": "Brazilian Portuguese",
		"ru":    "Russian",
		"ja":    "Japanese",
		"ko":    "Korean",
		"zh":    "Chinese (Simplified)",
		"zh-TW": "Chinese (Traditional)",
		"ar":    "Arabic",
		"hi":    "Hindi",
		"nl":    "Dutch",
		"pl":    "Polish",
		"sv":    "Swedish",
		"no":    "Norwegian",
		"da":    "Danish",
		"fi":    "Finnish",
		"el":    "Greek",
		"cs":    "Czech",
		"hu":    "Hungarian",
		"ro":    "Romanian",
		"bg":    "Bulgarian",
		"uk":    "Ukrainian",
		"hr":    "Croatian",
		"sr":    "Serbian",
		"sk":    "Slovak",
		"sl":    "Slovenian",
		"th":    "Thai",
		"vi":    "Vietnamese",
		"id":    "Indonesian",
		"ms":    "Malay",
		"tr":    "Turkish",
		"he":    "Hebrew",
		"fa":    "Persian",
		"ur":    "Urdu",
		"bn":    "Bengali",
		"ta":    "Tamil",
		"te":    "Telugu",
		"mr":    "Marathi",
		"gu":    "Gujarati",
		"kn":    "Kannada",
		"ml":    "Malayalam",
		"fil":   "Filipino",
		"sw":    "Swahili",
		"af":    "Afrikaans",
		"ca":    "Catalan",
		"eu":    "Basque",
		"ga":    "Irish",
		"cy":    "Welsh",
		"is":    "Icelandic",
		"et":    "Estonian",
		"lv":    "Latvian",
		"lt":    "Lithuanian",
		"mt":    "Maltese",
		"sq":    "Albanian",
		"mk":    "Macedonian",
		"auto":  "Auto-detect",
	}
	if name, ok := languages[code]; ok {
		return name
	}
	return code
}

func (a *TranslateTextAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[translateTextActionProps](ctx)
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

	if input.Text == "" {
		return nil, fmt.Errorf("text to translate cannot be empty")
	}

	if input.TargetLang == "" {
		return nil, fmt.Errorf("target language must be specified")
	}

	if input.Temperature == 0 {
		input.Temperature = 0.3
	}

	if input.SourceLang == "" {
		input.SourceLang = "auto"
	}

	if input.Formality == "" {
		input.Formality = "auto"
	}

	var prompt string

	systemPrompt := "You are a professional translator with native-level fluency in multiple languages. "
	systemPrompt += "Provide accurate, natural-sounding translations that preserve meaning, context, and cultural nuances. "
	systemPrompt += "When translating, consider idiomatic expressions and provide culturally appropriate equivalents rather than literal translations."

	if input.SourceLang == "auto" {
		prompt = fmt.Sprintf("Translate the following text to %s.", getLanguageName(input.TargetLang))
	} else {
		prompt = fmt.Sprintf("Translate the following %s text to %s.",
			getLanguageName(input.SourceLang), getLanguageName(input.TargetLang))
	}

	switch input.Formality {
	case "formal":
		prompt += " Use formal language and appropriate honorifics."
	case "informal":
		prompt += " Use informal, conversational language."
	case "business":
		prompt += " Use professional business language."
	case "casual":
		prompt += " Use casual, friendly language."
	case "academic":
		prompt += " Use academic language suitable for scholarly work."
	}

	if input.PreserveTone {
		prompt += " Preserve the original tone, style, and emotional content of the text."
	}

	prompt += fmt.Sprintf("\n\nText to translate:\n%s\n\n", input.Text)
	prompt += "Provide your response in the following JSON format:\n"
	prompt += `{
  "translated_text": "the main translation",
  "detected_language": "detected source language code if auto-detection was used",
  "literal_translation": "optional literal translation if significantly different from main translation",
  "alternative": "an alternative translation if applicable",
  "notes": "any cultural or contextual notes about the translation",
  "confidence": 0.95
}`

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
		System:      systemPrompt,
		Temperature: input.Temperature,
		MaxTokens:   4096,
	}

	gctx := context.Background()
	response, err := shared.CallClaudeAPI(gctx, authCtx, request)
	if err != nil {
		return nil, fmt.Errorf("translation failed: %w", err)
	}

	responseText := shared.ExtractResponseText(response)

	var translationResult map[string]interface{}
	if err := json.Unmarshal([]byte(responseText), &translationResult); err != nil {
		return map[string]interface{}{
			"translated_text": responseText,
			"source_language": getLanguageName(input.SourceLang),
			"target_language": getLanguageName(input.TargetLang),
			"model":           response.Model,
			"usage": map[string]int{
				"input_tokens":  response.Usage.InputTokens,
				"output_tokens": response.Usage.OutputTokens,
				"total_tokens":  response.Usage.InputTokens + response.Usage.OutputTokens,
			},
		}, nil
	}

	result := map[string]interface{}{
		"source_language": getLanguageName(input.SourceLang),
		"target_language": getLanguageName(input.TargetLang),
		"formality":       input.Formality,
		"model":           response.Model,
		"usage": map[string]int{
			"input_tokens":  response.Usage.InputTokens,
			"output_tokens": response.Usage.OutputTokens,
			"total_tokens":  response.Usage.InputTokens + response.Usage.OutputTokens,
		},
	}

	for key, value := range translationResult {
		result[key] = value
	}

	return result, nil
}

func NewTranslateTextAction() sdk.Action {
	return &TranslateTextAction{}
}
