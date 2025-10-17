// analyze_text.go - Text analysis action for Claude integration

package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/claude/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type analyzeTextActionProps struct {
	Text         string  `json:"text"`
	AnalysisType string  `json:"analysis_type"`
	Model        string  `json:"model"`
	DetailLevel  string  `json:"detail_level"`
	Language     string  `json:"language"`
	CustomFocus  string  `json:"custom_focus"`
	Temperature  float64 `json:"temperature"`
}

type AnalyzeTextAction struct{}

func (a *AnalyzeTextAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "analyze_text_claude",
		DisplayName:   "Analyze Text",
		Description:   "Perform deep text analysis including sentiment, tone, themes, entities, and style using Claude's advanced language understanding.",
		Type:          core.ActionTypeAction,
		Documentation: analyzeTextDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"analysis": map[string]any{
				"sentiment": map[string]any{
					"primary":    "positive",
					"confidence": 0.85,
					"scores": map[string]float64{
						"positive": 0.85,
						"negative": 0.10,
						"neutral":  0.05,
					},
					"explanation": "The text expresses optimism and enthusiasm",
				},
				"tone": map[string]any{
					"primary":   "professional",
					"secondary": []string{"confident", "informative"},
					"formality": "formal",
				},
				"themes": []map[string]any{
					{
						"theme":       "innovation",
						"relevance":   0.9,
						"description": "Focus on technological advancement",
					},
					{
						"theme":       "growth",
						"relevance":   0.7,
						"description": "Business expansion topics",
					},
				},
				"key_points": []string{
					"Main argument about technological advancement",
					"Supporting evidence from recent studies",
					"Conclusion about future implications",
				},
				"style": map[string]any{
					"complexity":         "moderate",
					"readability":        "college level",
					"vocabulary_level":   "advanced",
					"sentence_variation": "high",
				},
				"statistics": map[string]int{
					"word_count":          500,
					"sentence_count":      25,
					"paragraph_count":     5,
					"avg_sentence_length": 20,
				},
			},
			"analysis_type": "comprehensive",
			"detail_level":  "standard",
			"model":         "claude-3-5-sonnet-20241022",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *AnalyzeTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("analyze_text_claude", "Analyze Text")

	shared.RegisterModelProps(form)

	form.TextareaField("text", "Text to Analyze").
		Placeholder("Paste the text you want to analyze...").
		HelpText("The text content to analyze (up to 100,000 characters)").
		Required(true)

	form.SelectField("analysis_type", "Analysis Type").
		Placeholder("Select analysis type").
		Required(true).
		AddOptions(
			smartform.NewOption("sentiment", "Sentiment Analysis"),
			smartform.NewOption("tone", "Tone Detection"),
			smartform.NewOption("themes", "Theme Extraction"),
			smartform.NewOption("key_points", "Key Points"),
			smartform.NewOption("style", "Writing Style"),
			smartform.NewOption("entities", "Entity Recognition"),
			smartform.NewOption("readability", "Readability Analysis"),
			smartform.NewOption("emotions", "Emotional Analysis"),
			smartform.NewOption("bias", "Bias Detection"),
			smartform.NewOption("factuality", "Fact vs Opinion"),
			smartform.NewOption("comprehensive", "Comprehensive Analysis (All)"),
		).
		HelpText("Type of analysis to perform")

	form.SelectField("detail_level", "Detail Level").
		Placeholder("Select detail level").
		Required(false).
		AddOptions(
			smartform.NewOption("basic", "Basic - Quick overview"),
			smartform.NewOption("standard", "Standard - Balanced detail"),
			smartform.NewOption("detailed", "Detailed - In-depth analysis"),
			smartform.NewOption("expert", "Expert - Maximum detail"),
		).
		HelpText("Level of detail in the analysis")

	form.SelectField("language", "Text Language").
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
			smartform.NewOption("ja", "Japanese"),
			smartform.NewOption("ko", "Korean"),
			smartform.NewOption("zh", "Chinese"),
			smartform.NewOption("ar", "Arabic"),
			smartform.NewOption("ru", "Russian"),
			smartform.NewOption("hi", "Hindi"),
			smartform.NewOption("nl", "Dutch"),
			smartform.NewOption("pl", "Polish"),
			smartform.NewOption("tr", "Turkish"),
		).
		HelpText("Language of the text (for better analysis)")

	form.TextareaField("custom_focus", "Custom Focus").
		Placeholder("Specify any particular aspects you want to focus on...").
		HelpText("Optional: Specific aspects or criteria to focus on in the analysis").
		Required(false)

	form.NumberField("temperature", "Temperature").
		Placeholder("0.2").
		HelpText("Controls analysis creativity (0=consistent, 1=creative, default: 0.2)").
		Required(false)

	return form.Build()
}

func (a *AnalyzeTextAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *AnalyzeTextAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[analyzeTextActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if input.Model == "" {
		return nil, errors.New("Model is required")
	}

	if authCtx.Extra == nil {
		return nil, errors.New("please add your claude api key to continue")
	}

	if input.Text == "" {
		return nil, errors.New("text to analyze cannot be empty")
	}

	// Set defaults
	if input.Temperature == 0 {
		input.Temperature = 0.2 // Lower temperature for consistent analysis
	}

	if input.DetailLevel == "" {
		input.DetailLevel = "standard"
	}

	if input.Language == "" {
		input.Language = "auto"
	}

	// Build analysis prompt based on type
	var prompt string
	systemPrompt := "You are an expert text analyst with deep expertise in linguistics, psychology, and data analysis. Provide accurate, insightful analysis with specific examples from the text. Always return your analysis in valid JSON format."

	switch input.AnalysisType {
	case "sentiment":
		prompt = buildSentimentPrompt(input.DetailLevel)
	case "tone":
		prompt = buildTonePrompt(input.DetailLevel)
	case "themes":
		prompt = buildThemesPrompt(input.DetailLevel)
	case "key_points":
		prompt = buildKeyPointsPrompt(input.DetailLevel)
	case "style":
		prompt = buildStylePrompt(input.DetailLevel)
	case "entities":
		prompt = buildEntitiesPrompt(input.DetailLevel)
	case "readability":
		prompt = buildReadabilityPrompt(input.DetailLevel)
	case "emotions":
		prompt = buildEmotionsPrompt(input.DetailLevel)
	case "bias":
		prompt = buildBiasPrompt(input.DetailLevel)
	case "factuality":
		prompt = buildFactualityPrompt(input.DetailLevel)
	case "comprehensive":
		prompt = buildComprehensivePrompt(input.DetailLevel)
	default:
		prompt = buildComprehensivePrompt(input.DetailLevel)
	}

	if input.CustomFocus != "" {
		prompt += fmt.Sprintf("\n\nAdditional focus areas: %s", input.CustomFocus)
	}

	if input.Language != "auto" {
		prompt += fmt.Sprintf("\n\nNote: The text is in %s.", getLanguageNameList(input.Language))
	}

	prompt += fmt.Sprintf("\n\nText to analyze:\n%s\n\nProvide your analysis in valid JSON format only, with no additional text or markdown formatting.", input.Text)

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
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	responseText := shared.ExtractResponseText(response)

	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(responseText), &analysis); err != nil {
		return map[string]interface{}{
			"analysis": map[string]interface{}{
				"raw_response": responseText,
				"parse_error":  err.Error(),
				"note":         "The analysis was completed but couldn't be parsed as JSON",
			},
			"analysis_type": input.AnalysisType,
			"detail_level":  input.DetailLevel,
			"model":         response.Model,
			"usage": map[string]int{
				"input_tokens":  response.Usage.InputTokens,
				"output_tokens": response.Usage.OutputTokens,
				"total_tokens":  response.Usage.InputTokens + response.Usage.OutputTokens,
			},
		}, nil
	}

	// Add text statistics
	textStats := getBasicTextStats(input.Text)

	return map[string]interface{}{
		"analysis":      analysis,
		"analysis_type": input.AnalysisType,
		"detail_level":  input.DetailLevel,
		"text_stats":    textStats,
		"model":         response.Model,
		"usage": map[string]int{
			"input_tokens":  response.Usage.InputTokens,
			"output_tokens": response.Usage.OutputTokens,
			"total_tokens":  response.Usage.InputTokens + response.Usage.OutputTokens,
		},
	}, nil
}

// Helper function to get basic text statistics
func getBasicTextStats(text string) map[string]int {
	words := strings.Fields(text)
	sentences := strings.Count(text, ".") + strings.Count(text, "!") + strings.Count(text, "?")
	if sentences == 0 {
		sentences = 1
	}
	paragraphs := strings.Count(text, "\n\n") + 1

	avgSentenceLength := len(words) / sentences
	if sentences > 0 {
		avgSentenceLength = len(words) / sentences
	}

	return map[string]int{
		"characters":          len(text),
		"words":               len(words),
		"sentences":           sentences,
		"paragraphs":          paragraphs,
		"avg_sentence_length": avgSentenceLength,
	}
}

func getLanguageNameList(code string) string {
	languages := map[string]string{
		"en": "English", "es": "Spanish", "fr": "French",
		"de": "German", "it": "Italian", "pt": "Portuguese",
		"ja": "Japanese", "ko": "Korean", "zh": "Chinese",
		"ar": "Arabic", "ru": "Russian", "hi": "Hindi",
		"nl": "Dutch", "pl": "Polish", "tr": "Turkish",
		"auto": "Auto-detect",
	}
	if name, ok := languages[code]; ok {
		return name
	}
	return code
}

// Prompt building functions for each analysis type

func buildSentimentPrompt(detailLevel string) string {
	base := "Analyze the sentiment of the following text. "
	switch detailLevel {
	case "basic":
		return base + `Return a JSON object with:
{
  "sentiment": "positive/negative/neutral/mixed",
  "confidence": 0.0-1.0,
  "summary": "one-line summary"
}`
	case "detailed", "expert":
		return base + `Return a detailed JSON object with:
{
  "primary_sentiment": "positive/negative/neutral/mixed",
  "confidence": 0.0-1.0,
  "sentiment_scores": {
    "positive": 0.0-1.0,
    "negative": 0.0-1.0,
    "neutral": 0.0-1.0
  },
  "emotional_valence": -1.0 to 1.0,
  "subjectivity": 0.0-1.0,
  "sentiment_progression": "description of how sentiment changes through the text",
  "notable_phrases": ["phrases that strongly indicate sentiment"],
  "contradictions": ["any conflicting sentiments found"],
  "context_factors": ["factors that might affect interpretation"],
  "confidence_factors": ["what makes you confident/uncertain about this analysis"]
}`
	default:
		return base + `Return a JSON object with:
{
  "sentiment": "positive/negative/neutral/mixed",
  "confidence": 0.0-1.0,
  "sentiment_scores": {
    "positive": 0.0-1.0,
    "negative": 0.0-1.0,
    "neutral": 0.0-1.0
  },
  "explanation": "brief explanation of the sentiment",
  "key_indicators": ["main phrases that indicate the sentiment"]
}`
	}
}

func buildTonePrompt(detailLevel string) string {
	base := "Analyze the tone of the following text. "
	switch detailLevel {
	case "basic":
		return base + `Return a JSON object with:
{
  "primary_tone": "professional/casual/formal/friendly/etc",
  "description": "brief description of the tone"
}`
	case "detailed", "expert":
		return base + `Return a detailed JSON object with:
{
  "primary_tone": "main tone",
  "secondary_tones": ["additional detected tones"],
  "tone_strength": 0.0-1.0,
  "formality_level": "very formal/formal/neutral/informal/very informal",
  "emotional_undertones": ["underlying emotions detected"],
  "author_attitude": "description of author's attitude toward the subject",
  "audience_orientation": "how the text addresses its intended audience",
  "linguistic_markers": ["specific words/phrases that indicate the tone"],
  "tone_consistency": 0.0-1.0,
  "tone_shifts": ["any notable changes in tone and where they occur"],
  "cultural_context": "any cultural factors affecting tone interpretation"
}`
	default:
		return base + `Return a JSON object with:
{
  "primary_tone": "main tone",
  "secondary_tones": ["additional tones"],
  "formality": "formal/neutral/informal",
  "emotional_undertones": ["detected emotions"],
  "explanation": "brief explanation of the tone"
}`
	}
}

func buildThemesPrompt(detailLevel string) string {
	base := "Extract and analyze the main themes from the following text. "
	switch detailLevel {
	case "basic":
		return base + `Return a JSON object with:
{
  "themes": ["theme1", "theme2", "theme3"],
  "main_topic": "primary subject matter"
}`
	case "detailed", "expert":
		return base + `Return a detailed JSON object with:
{
  "main_themes": [
    {
      "theme": "theme name",
      "relevance": 0.0-1.0,
      "frequency": "how often it appears",
      "supporting_evidence": ["quotes or specific references"],
      "subtopics": ["related subtopics within this theme"]
    }
  ],
  "secondary_themes": ["less prominent themes"],
  "emerging_themes": ["themes that are developing but not fully formed"],
  "implicit_themes": ["underlying themes not explicitly stated"],
  "theme_relationships": "description of how themes connect and interact",
  "thematic_progression": "how themes develop throughout the text",
  "cultural_context": "cultural or contextual themes present",
  "symbolic_elements": ["any symbolic or metaphorical themes"]
}`
	default:
		return base + `Return a JSON object with:
{
  "themes": [
    {
      "theme": "theme name",
      "relevance": 0.0-1.0,
      "description": "brief description"
    }
  ],
  "main_theme": "most prominent theme",
  "theme_connections": "how themes relate to each other"
}`
	}
}

func buildKeyPointsPrompt(detailLevel string) string {
	base := "Extract the key points from the following text. "
	switch detailLevel {
	case "basic":
		return base + `Return a JSON object with:
{
  "key_points": ["point 1", "point 2", "point 3"],
  "main_idea": "central message"
}`
	case "detailed", "expert":
		return base + `Return a detailed JSON object with:
{
  "main_arguments": [
    {
      "point": "the key point",
      "importance": "high/medium/low",
      "supporting_details": ["evidence or examples"],
      "location": "where in text (beginning/middle/end)",
      "strength": "how well supported (strong/moderate/weak)"
    }
  ],
  "supporting_points": ["secondary points that support main arguments"],
  "conclusions": ["main conclusions drawn"],
  "recommendations": ["any recommendations or calls to action"],
  "assumptions": ["underlying assumptions made"],
  "implications": ["broader implications of the points"],
  "counter_arguments": ["any counter-arguments presented"],
  "evidence_quality": "assessment of evidence quality"
}`
	default:
		return base + `Return a JSON object with:
{
  "key_points": ["main points in order of importance"],
  "supporting_points": ["secondary supporting points"],
  "main_conclusion": "primary conclusion if any",
  "summary": "brief summary of main ideas"
}`
	}
}

func buildStylePrompt(detailLevel string) string {
	base := "Analyze the writing style of the following text. "
	switch detailLevel {
	case "basic":
		return base + `Return a JSON object with:
{
  "style": "academic/conversational/technical/narrative/etc",
  "complexity": "simple/moderate/complex",
  "readability": "easy/moderate/difficult"
}`
	case "detailed", "expert":
		return base + `Return a detailed JSON object with:
{
  "overall_style": "classification of writing style",
  "complexity_metrics": {
    "vocabulary_level": "basic/intermediate/advanced/expert",
    "sentence_complexity": "simple/compound/complex/varied",
    "average_sentence_length": estimated number,
    "paragraph_structure": "short/medium/long/varied"
  },
  "readability_assessment": {
    "grade_level": "estimated education level required",
    "clarity": 0.0-1.0,
    "accessibility": "very accessible/accessible/moderate/challenging"
  },
  "linguistic_features": {
    "voice": "active/passive/mixed",
    "tense": "primary tense used",
    "person": "first/second/third/mixed",
    "sentence_variety": "high/medium/low"
  },
  "stylistic_devices": ["metaphors", "analogies", "alliteration", "etc"],
  "vocabulary_characteristics": {
    "formality": "formal/neutral/informal",
    "technicality": "non-technical/somewhat technical/highly technical",
    "jargon_usage": "none/minimal/moderate/heavy"
  },
  "strengths": ["writing strengths identified"],
  "weaknesses": ["areas for improvement"],
  "distinctive_features": ["unique stylistic elements"]
}`
	default:
		return base + `Return a JSON object with:
{
  "style": "writing style classification",
  "complexity": "simple/moderate/complex",
  "readability": "grade level or difficulty",
  "vocabulary": "basic/intermediate/advanced",
  "sentence_structure": "simple/varied/complex",
  "notable_features": ["distinctive style characteristics"],
  "strengths": ["key strengths"],
  "improvements": ["suggested improvements"]
}`
	}
}

func buildEntitiesPrompt(detailLevel string) string {
	base := "Extract and categorize all named entities from the text. "
	if detailLevel == "detailed" || detailLevel == "expert" {
		return base + `Return a detailed JSON object with:
{
  "people": [
    {"name": "person name", "context": "role or description", "mentions": count}
  ],
  "organizations": [
    {"name": "org name", "type": "company/institution/etc", "mentions": count}
  ],
  "locations": [
    {"name": "place name", "type": "city/country/etc", "mentions": count}
  ],
  "dates": [
    {"date": "date reference", "context": "what it refers to"}
  ],
  "numbers": [
    {"value": "number", "context": "what it represents"}
  ],
  "products": ["product or service names"],
  "events": ["event names"],
  "technical_terms": ["domain-specific terminology"],
  "acronyms": [
    {"acronym": "ABC", "expansion": "full form if known"}
  ],
  "other_entities": ["other notable entities"]
}`
	}
	return base + `Return a JSON object with:
{
  "people": ["person names"],
  "organizations": ["company/organization names"],
  "locations": ["place names"],
  "dates": ["temporal references"],
  "numbers": ["significant quantities or amounts"],
  "products": ["product or service names"],
  "technical_terms": ["specialized terminology"],
  "other": ["other notable entities"]
}`
}

func buildReadabilityPrompt(detailLevel string) string {
	return `Analyze the readability of the text. Return a JSON object with:
{
  "overall_readability": "very easy/easy/moderate/difficult/very difficult",
  "target_audience": "description of intended readers",
  "education_level": "estimated grade level required",
  "accessibility_score": 0.0-1.0,
  "clarity_score": 0.0-1.0,
  "complexity_factors": ["factors making it complex"],
  "jargon_assessment": {
    "level": "none/minimal/moderate/heavy",
    "examples": ["jargon terms found"]
  },
  "sentence_analysis": {
    "variety": "good/moderate/poor",
    "average_length": "short/medium/long"
  },
  "improvement_suggestions": ["specific ways to improve readability"],
  "strengths": ["what makes it readable"],
  "challenges": ["what makes it difficult"]
}`
}

func buildEmotionsPrompt(detailLevel string) string {
	return `Analyze the emotions expressed in the text. Return a JSON object with:
{
  "primary_emotion": "dominant emotion detected",
  "emotion_spectrum": {
    "joy": 0.0-1.0,
    "sadness": 0.0-1.0,
    "anger": 0.0-1.0,
    "fear": 0.0-1.0,
    "surprise": 0.0-1.0,
    "disgust": 0.0-1.0,
    "trust": 0.0-1.0,
    "anticipation": 0.0-1.0
  },
  "emotional_intensity": 0.0-1.0,
  "emotional_stability": "stable/variable/volatile",
  "emotional_progression": "how emotions change through the text",
  "trigger_phrases": ["phrases with strong emotional content"],
  "emotional_appeals": ["any emotional appeals made"],
  "author_emotion": "apparent emotional state of the author",
  "intended_emotion": "emotions the text aims to evoke in readers"
}`
}

func buildBiasPrompt(detailLevel string) string {
	return `Analyze potential biases in the text. Return a JSON object with:
{
  "overall_bias_level": "none/minimal/moderate/significant/heavy",
  "detected_biases": [
    {
      "type": "political/cultural/gender/racial/confirmation/etc",
      "evidence": ["specific examples from text"],
      "severity": "low/medium/high"
    }
  ],
  "perspective": "description of author's apparent perspective",
  "balance_assessment": "balanced/somewhat biased/heavily biased",
  "missing_perspectives": ["viewpoints not represented"],
  "loaded_language": ["biased or charged terms used"],
  "assumptions": ["unstated assumptions present"],
  "objectivity_score": 0.0-1.0,
  "recommendations": ["suggestions for more balanced presentation"],
  "neutral_alternatives": ["neutral alternatives for biased language"]
}`
}

func buildFactualityPrompt(detailLevel string) string {
	return `Analyze the factuality and opinion content of the text. Return a JSON object with:
{
  "fact_opinion_ratio": "mostly facts/balanced/mostly opinions",
  "factual_claims": [
    {"claim": "statement", "verifiability": "easily verifiable/difficult to verify/unverifiable"}
  ],
  "opinions": [
    {"opinion": "statement", "attribution": "who expressed it"}
  ],
  "unsupported_claims": ["claims lacking evidence"],
  "evidence_quality": "strong/moderate/weak/absent",
  "sources_cited": ["any sources referenced"],
  "credibility_indicators": ["what makes it credible"],
  "red_flags": ["what raises credibility concerns"],
  "verification_needed": ["claims that should be fact-checked"],
  "logical_fallacies": ["any logical fallacies detected"],
  "overall_credibility": 0.0-1.0
}`
}

func buildComprehensivePrompt(detailLevel string) string {
	if detailLevel == "basic" {
		return `Perform a comprehensive but concise analysis of the text. Return a JSON object with:
{
  "sentiment": {
    "primary": "positive/negative/neutral/mixed",
    "confidence": 0.0-1.0
  },
  "tone": {
    "primary": "main tone",
    "formality": "formal/neutral/informal"
  },
  "themes": ["top 3-5 themes"],
  "key_points": ["main 3-5 points"],
  "style": {
    "type": "writing style",
    "complexity": "simple/moderate/complex"
  },
  "entities": {
    "people": [],
    "organizations": [],
    "locations": []
  },
  "summary": "2-3 sentence summary"
}`
	}

	return `Perform a comprehensive analysis of the text including all aspects. Return a detailed JSON object with:
{
  "sentiment": {
    "primary": "positive/negative/neutral/mixed",
    "confidence": 0.0-1.0,
    "scores": {"positive": 0.0, "negative": 0.0, "neutral": 0.0},
    "explanation": "sentiment explanation"
  },
  "tone": {
    "primary": "main tone",
    "secondary": ["additional tones"],
    "formality": "formality level",
    "consistency": 0.0-1.0
  },
  "themes": [
    {"theme": "theme name", "relevance": 0.0-1.0, "description": "brief description"}
  ],
  "key_points": ["main points extracted"],
  "style": {
    "type": "writing style",
    "complexity": "complexity level",
    "readability": "readability level",
    "vocabulary": "vocabulary level",
    "strengths": ["style strengths"],
    "weaknesses": ["style weaknesses"]
  },
  "entities": {
    "people": ["person names"],
    "organizations": ["org names"],
    "locations": ["place names"],
    "dates": ["date references"],
    "technical_terms": ["specialized terms"]
  },
  "emotions": {
    "primary": "dominant emotion",
    "intensity": 0.0-1.0,
    "spectrum": {"joy": 0.0, "sadness": 0.0, "anger": 0.0, "fear": 0.0}
  },
  "bias_assessment": {
    "level": "none/minimal/moderate/significant",
    "types": ["detected bias types"],
    "objectivity": 0.0-1.0
  },
  "factuality": {
    "fact_opinion_ratio": "mostly facts/balanced/mostly opinions",
    "credibility": 0.0-1.0,
    "evidence_quality": "strong/moderate/weak"
  },
  "statistics": {
    "word_count": number,
    "sentence_count": number,
    "paragraph_count": number,
    "avg_sentence_length": number
  },
  "summary": "comprehensive summary of the text",
  "recommendations": ["suggestions for improvement if applicable"]
}`
}

func NewAnalyzeTextAction() sdk.Action {
	return &AnalyzeTextAction{}
}
