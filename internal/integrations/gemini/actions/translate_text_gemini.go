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

type translateTextActionProps struct {
	Text       string `json:"text"`
	TargetLang string `json:"target_lang"`
	SourceLang string `json:"source_lang"`
	Model      string `json:"model"`
}

type TranslateTextAction struct{}

func (a *TranslateTextAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "translate_text_gemini",
		DisplayName:   "Translate Text",
		Description:   "Translate text between multiple languages with high accuracy using Gemini's multilingual capabilities.",
		Type:          core.ActionTypeAction,
		Documentation: translateTextGeminiDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"translated_text": "Bonjour le monde!",
			"source_language": "en",
			"target_language": "fr",
			"model":           "gemini-1.5-flash",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *TranslateTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("translate_text_gemini", "Translate Text")

	form.TextareaField("text", "Text to Translate").
		Placeholder("Enter text to translate").
		HelpText("The text you want to translate").
		Required(true)

	form.SelectField("target_lang", "Target Language").
		Placeholder("Select target language").
		Required(true).
		AddOptions(
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
			smartform.NewOption("ms", "Malay"),
		).
		HelpText("Language to translate to")

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
			smartform.NewOption("zh", "Chinese"),
			smartform.NewOption("ar", "Arabic"),
			smartform.NewOption("hi", "Hindi"),
		).
		HelpText("Language to translate from (leave empty for auto-detect)")

	RegisterModelProps(form)

	return form.Build()
}

func (a *TranslateTextAction) Auth() *core.AuthMetadata {
	return nil
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

	gctx := context.Background()
	client, err := CreateGeminiClient(gctx, authCtx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	modelName := strings.TrimPrefix(input.Model, "models/")
	model := client.GenerativeModel(modelName)

	// Build translation prompt
	var prompt string
	if input.SourceLang == "" || input.SourceLang == "auto" {
		prompt = fmt.Sprintf("Translate the following text to %s. Preserve the tone and style:\n\n%s",
			getLanguageName(input.TargetLang), input.Text)
	} else {
		prompt = fmt.Sprintf("Translate the following %s text to %s. Preserve the tone and style:\n\n%s",
			getLanguageName(input.SourceLang), getLanguageName(input.TargetLang), input.Text)
	}

	content, err := model.GenerateContent(gctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	var translatedText string
	if content != nil && len(content.Candidates) > 0 {
		for _, part := range content.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				translatedText += string(text)
			}
		}
	}

	return map[string]interface{}{
		"translated_text": translatedText,
		"source_language": input.SourceLang,
		"target_language": input.TargetLang,
		"model":           modelName,
	}, nil
}

func getLanguageName(code string) string {
	languages := map[string]string{
		"en": "English", "es": "Spanish", "fr": "French",
		"de": "German", "it": "Italian", "pt": "Portuguese",
		"ru": "Russian", "ja": "Japanese", "ko": "Korean",
		"zh": "Chinese", "ar": "Arabic", "hi": "Hindi",
		"nl": "Dutch", "pl": "Polish", "tr": "Turkish",
		"vi": "Vietnamese", "th": "Thai", "id": "Indonesian",
		"ms": "Malay", "zh-TW": "Traditional Chinese",
	}
	if name, ok := languages[code]; ok {
		return name
	}
	return code
}

func NewTranslateTextAction() sdk.Action {
	return &TranslateTextAction{}
}
