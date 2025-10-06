// compare_texts.go - Text comparison action for Claude integration

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

type compareTextsActionProps struct {
	Text1            string  `json:"text1"`
	Text2            string  `json:"text2"`
	ComparisonType   string  `json:"comparison_type"`
	Model            string  `json:"model"`
	DetailLevel      string  `json:"detail_level"`
	IgnoreCase       bool    `json:"ignore_case"`
	IgnoreWhitespace bool    `json:"ignore_whitespace"`
	FocusAreas       string  `json:"focus_areas"`
	Temperature      float64 `json:"temperature"`
}

type CompareTextsAction struct{}

func (a *CompareTextsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "compare_texts_claude",
		DisplayName:   "Compare Texts",
		Description:   "Compare two texts for similarities, differences, style, content, and more using Claude's advanced analysis capabilities.",
		Type:          core.ActionTypeAction,
		Documentation: compareTextsDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"comparison": map[string]any{
				"overall_similarity": 0.75,
				"similarity_breakdown": map[string]float64{
					"content":   0.80,
					"style":     0.70,
					"structure": 0.75,
					"semantic":  0.85,
				},
				"differences": map[string]any{
					"major": []string{"Text 1 focuses on X while Text 2 emphasizes Y"},
					"minor": []string{"Different examples used"},
					"style": []string{"Text 1 is more formal"},
				},
				"similarities": []string{
					"Both texts discuss the same main topic",
					"Similar conclusion reached",
				},
				"unique_to_text1": []string{"Specific point about A"},
				"unique_to_text2": []string{"Additional detail about B"},
				"statistics": map[string]map[string]int{
					"text1": {"words": 500, "sentences": 25},
					"text2": {"words": 450, "sentences": 22},
				},
				"recommendation": "The texts are substantially similar with minor variations",
			},
			"comparison_type": "comprehensive",
			"model":           "claude-3-5-sonnet-20241022",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CompareTextsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("compare_texts_claude", "Compare Texts")

	shared.RegisterModelProps(form)

	form.TextareaField("text1", "First Text").
		Placeholder("Paste the first text here...").
		HelpText("The first text to compare").
		Required(true)

	form.TextareaField("text2", "Second Text").
		Placeholder("Paste the second text here...").
		HelpText("The second text to compare").
		Required(true)

	form.SelectField("comparison_type", "Comparison Type").
		Placeholder("Select comparison type").
		Required(true).
		AddOptions(
			smartform.NewOption("similarity", "Similarity Analysis"),
			smartform.NewOption("differences", "Difference Detection"),
			smartform.NewOption("style", "Style Comparison"),
			smartform.NewOption("content", "Content Comparison"),
			smartform.NewOption("plagiarism", "Plagiarism Check"),
			smartform.NewOption("version", "Version Comparison"),
			smartform.NewOption("translation", "Translation Verification"),
			smartform.NewOption("quality", "Quality Comparison"),
			smartform.NewOption("factual", "Factual Consistency"),
			smartform.NewOption("comprehensive", "Comprehensive Comparison"),
		).
		HelpText("Type of comparison to perform")

	form.SelectField("detail_level", "Detail Level").
		Placeholder("Select detail level").
		Required(false).
		AddOptions(
			smartform.NewOption("basic", "Basic - Quick overview"),
			smartform.NewOption("standard", "Standard - Balanced detail"),
			smartform.NewOption("detailed", "Detailed - In-depth analysis"),
			smartform.NewOption("expert", "Expert - Maximum detail"),
		).
		HelpText("Level of detail in the comparison")

	form.CheckboxField("ignore_case", "Ignore Case").
		HelpText("Ignore letter case differences in comparison").
		Required(false)

	form.CheckboxField("ignore_whitespace", "Ignore Whitespace").
		HelpText("Ignore whitespace and formatting differences").
		Required(false)

	form.TextareaField("focus_areas", "Focus Areas").
		Placeholder("Specify particular aspects to focus on...").
		HelpText("Optional: Specific aspects or criteria to emphasize in the comparison").
		Required(false)

	form.NumberField("temperature", "Temperature").
		Placeholder("0.2").
		HelpText("Controls analysis creativity (0=consistent, 1=creative, default: 0.2)").
		Required(false)

	return form.Build()
}

func (a *CompareTextsAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *CompareTextsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[compareTextsActionProps](ctx)
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

	if input.Text1 == "" || input.Text2 == "" {
		return nil, errors.New("both texts are required for comparison")
	}

	if input.Temperature == 0 {
		input.Temperature = 0.2 
	}

	if input.DetailLevel == "" {
		input.DetailLevel = "standard"
	}

	var prompt string
	systemPrompt := "You are an expert text analyst specializing in comparative analysis. Provide detailed, objective comparisons with specific examples and quantifiable metrics where possible."

	preprocessInstructions := ""
	if input.IgnoreCase {
		preprocessInstructions += "\nNote: Ignore case differences when comparing."
	}
	if input.IgnoreWhitespace {
		preprocessInstructions += "\nNote: Ignore whitespace and formatting differences when comparing."
	}

	switch input.ComparisonType {
	case "similarity":
		prompt = buildSimilarityComparisonPrompt(input.DetailLevel)
	case "differences":
		prompt = buildDifferenceComparisonPrompt(input.DetailLevel)
	case "style":
		prompt = buildStyleComparisonPrompt(input.DetailLevel)
	case "content":
		prompt = buildContentComparisonPrompt(input.DetailLevel)
	case "plagiarism":
		prompt = buildPlagiarismComparisonPrompt(input.DetailLevel)
	case "version":
		prompt = buildVersionComparisonPrompt(input.DetailLevel)
	case "translation":
		prompt = buildTranslationComparisonPrompt(input.DetailLevel)
	case "quality":
		prompt = buildQualityComparisonPrompt(input.DetailLevel)
	case "factual":
		prompt = buildFactualComparisonPrompt(input.DetailLevel)
	case "comprehensive":
		prompt = buildComprehensiveComparisonPrompt(input.DetailLevel)
	default:
		prompt = buildComprehensiveComparisonPrompt(input.DetailLevel)
	}

	prompt += preprocessInstructions

	if input.FocusAreas != "" {
		prompt += fmt.Sprintf("\n\nFocus particularly on: %s", input.FocusAreas)
	}

	prompt += fmt.Sprintf("\n\nText 1:\n%s\n\nText 2:\n%s\n\nProvide your comparison in JSON format.", input.Text1, input.Text2)

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
		return nil, fmt.Errorf("comparison failed: %w", err)
	}

	responseText := shared.ExtractResponseText(response)

	var comparison map[string]interface{}
	if err := json.Unmarshal([]byte(responseText), &comparison); err != nil {
		comparison = map[string]interface{}{
			"raw_comparison": responseText,
			"parse_error":    err.Error(),
		}
	}

	text1Stats := getTextStatistics(input.Text1)
	text2Stats := getTextStatistics(input.Text2)

	return map[string]interface{}{
		"comparison":      comparison,
		"comparison_type": input.ComparisonType,
		"detail_level":    input.DetailLevel,
		"text_statistics": map[string]interface{}{
			"text1": text1Stats,
			"text2": text2Stats,
		},
		"model": response.Model,
		"usage": map[string]int{
			"input_tokens":  response.Usage.InputTokens,
			"output_tokens": response.Usage.OutputTokens,
			"total_tokens":  response.Usage.InputTokens + response.Usage.OutputTokens,
		},
	}, nil
}

func getTextStatistics(text string) map[string]int {
	words := len(strings.Fields(text))
	sentences := len(strings.Split(text, "."))
	paragraphs := len(strings.Split(text, "\n\n"))
	characters := len(text)

	return map[string]int{
		"words":      words,
		"sentences":  sentences,
		"paragraphs": paragraphs,
		"characters": characters,
	}
}


func buildSimilarityComparisonPrompt(detailLevel string) string {
	base := "Compare these two texts and analyze their similarity. Return a JSON object with:\n"
	switch detailLevel {
	case "basic":
		return base + `{
  "similarity_score": 0.0-1.0,
  "verdict": "very similar/somewhat similar/different",
  "main_similarities": ["key similarities"]
}`
	case "detailed", "expert":
		return base + `{
  "overall_similarity": 0.0-1.0,
  "similarity_breakdown": {
    "content_similarity": 0.0-1.0,
    "style_similarity": 0.0-1.0,
    "structure_similarity": 0.0-1.0,
    "semantic_similarity": 0.0-1.0,
    "lexical_overlap": 0.0-1.0,
    "conceptual_alignment": 0.0-1.0
  },
  "shared_elements": {
    "themes": ["common themes"],
    "key_points": ["shared key points"],
    "vocabulary": ["common significant terms"],
    "phrases": ["identical or near-identical phrases"],
    "arguments": ["shared arguments or reasoning"]
  },
  "similarity_matrix": "detailed breakdown of paragraph-by-paragraph similarity",
  "statistical_measures": {
    "jaccard_index": "estimate",
    "cosine_similarity": "estimate",
    "edit_distance": "normalized estimate"
  },
  "verdict": "detailed assessment",
  "confidence": 0.0-1.0
}`
	default:
		return base + `{
  "similarity_score": 0.0-1.0,
  "content_similarity": 0.0-1.0,
  "style_similarity": 0.0-1.0,
  "main_similarities": ["key similarities"],
  "shared_themes": ["common themes"],
  "verdict": "assessment of similarity"
}`
	}
}

func buildDifferenceComparisonPrompt(detailLevel string) string {
	base := "Compare these two texts and identify all differences. Return a JSON object with:\n"
	switch detailLevel {
	case "basic":
		return base + `{
  "major_differences": ["significant differences"],
  "minor_differences": ["small differences"],
  "verdict": "summary of differences"
}`
	case "detailed", "expert":
		return base + `{
  "content_differences": {
    "major": ["significant content differences"],
    "minor": ["minor content variations"],
    "factual": ["factual discrepancies"],
    "examples": ["different examples or evidence used"]
  },
  "structural_differences": {
    "organization": "how structure differs",
    "length": "length comparison",
    "sections": ["sections in one but not other"],
    "flow": "difference in logical flow"
  },
  "stylistic_differences": {
    "tone": "tone variations",
    "formality": "formality differences",
    "voice": "active vs passive",
    "complexity": "complexity differences"
  },
  "unique_to_text1": ["content only in text 1"],
  "unique_to_text2": ["content only in text 2"],
  "contradictions": ["direct contradictions between texts"],
  "difference_significance": "assessment of how significant differences are",
  "difference_patterns": "patterns in the differences"
}`
	default:
		return base + `{
  "major_differences": ["significant differences"],
  "minor_differences": ["small differences"],
  "unique_to_text1": ["content only in text 1"],
  "unique_to_text2": ["content only in text 2"],
  "style_differences": ["style variations"],
  "verdict": "summary of differences"
}`
	}
}

func buildStyleComparisonPrompt(detailLevel string) string {
	base := "Compare the writing styles of these two texts. Return a JSON object with:\n"
	switch detailLevel {
	case "basic":
		return base + `{
  "text1_style": "style description",
  "text2_style": "style description",
  "main_differences": ["key style differences"]
}`
	default:
		return base + `{
  "text1_style": {
    "overall": "style classification",
    "tone": "tone description",
    "formality": "formality level",
    "complexity": "complexity level"
  },
  "text2_style": {
    "overall": "style classification",
    "tone": "tone description",
    "formality": "formality level",
    "complexity": "complexity level"
  },
  "comparative_analysis": {
    "tone_comparison": "how tones compare",
    "formality_comparison": "formality differences",
    "vocabulary_comparison": "vocabulary sophistication",
    "sentence_structure": "sentence complexity comparison",
    "paragraph_structure": "organization comparison"
  },
  "stylistic_similarities": ["shared style features"],
  "stylistic_differences": ["contrasting style features"],
  "readability_comparison": "which is easier to read and why",
  "audience_suitability": "different audiences each style suits"
}`
	}
}

func buildContentComparisonPrompt(detailLevel string) string {
	return `Compare the content and topics of these two texts. Return a JSON object with:
{
  "topic_overlap": {
    "shared_topics": ["topics in both texts"],
    "overlap_percentage": 0.0-1.0
  },
  "unique_content": {
    "text1_only": ["topics only in text 1"],
    "text2_only": ["topics only in text 2"]
  },
  "depth_comparison": {
    "text1_depth": "how deeply text 1 covers topics",
    "text2_depth": "how deeply text 2 covers topics"
  },
  "argument_comparison": {
    "shared_arguments": ["common arguments"],
    "conflicting_arguments": ["contradictory arguments"],
    "complementary_arguments": ["arguments that complement each other"]
  },
  "evidence_comparison": "comparison of evidence/examples used",
  "conclusion_comparison": "how conclusions compare",
  "content_quality": "relative quality assessment"
}`
}

func buildPlagiarismComparisonPrompt(detailLevel string) string {
	return `Check for potential plagiarism between these texts. Return a JSON object with:
{
  "plagiarism_score": 0.0-1.0,
  "identical_phrases": ["exact matches if any"],
  "paraphrased_sections": ["likely paraphrased content"],
  "structural_copying": "similar structure patterns",
  "idea_similarity": "similarity of ideas vs expression",
  "originality_assessment": {
    "text1_originality": 0.0-1.0,
    "text2_originality": 0.0-1.0
  },
  "potential_source": "which might be source if copying detected",
  "verdict": "plagiarism assessment",
  "confidence": 0.0-1.0,
  "recommendation": "suggested action"
}`
}

func buildVersionComparisonPrompt(detailLevel string) string {
	return `Compare these texts as potential versions of the same document. Return a JSON object with:
{
  "version_relationship": "likely same document/different documents",
  "changes": {
    "additions": ["content added"],
    "deletions": ["content removed"],
    "modifications": ["content changed"],
    "reorganizations": ["structural changes"]
  },
  "change_statistics": {
    "added_words": number,
    "deleted_words": number,
    "modified_sentences": number
  },
  "improvement_assessment": {
    "clarity_change": "improved/degraded/unchanged",
    "completeness_change": "more/less complete",
    "quality_change": "better/worse/similar"
  },
  "version_recommendation": "which version is preferable and why",
  "change_summary": "overall summary of changes"
}`
}

func buildTranslationComparisonPrompt(detailLevel string) string {
	return `Compare these texts as potential translations or multilingual versions. Return a JSON object with:
{
  "translation_likelihood": 0.0-1.0,
  "semantic_equivalence": 0.0-1.0,
  "detected_languages": {
    "text1": "detected or specified language",
    "text2": "detected or specified language"
  },
  "translation_quality": {
    "accuracy": 0.0-1.0,
    "fluency": 0.0-1.0,
    "completeness": 0.0-1.0
  },
  "translation_issues": [
    {"issue": "description", "severity": "minor/major"}
  ],
  "cultural_adaptation": "how well cultural elements are handled",
  "missing_content": ["content not translated"],
  "added_content": ["content added in translation"],
  "recommendation": "translation quality assessment"
}`
}

func buildQualityComparisonPrompt(detailLevel string) string {
	return `Compare the quality of these two texts. Return a JSON object with:
{
  "quality_scores": {
    "text1": 0.0-1.0,
    "text2": 0.0-1.0
  },
  "quality_dimensions": {
    "clarity": {"text1": 0.0-1.0, "text2": 0.0-1.0},
    "coherence": {"text1": 0.0-1.0, "text2": 0.0-1.0},
    "completeness": {"text1": 0.0-1.0, "text2": 0.0-1.0},
    "accuracy": {"text1": 0.0-1.0, "text2": 0.0-1.0},
    "engagement": {"text1": 0.0-1.0, "text2": 0.0-1.0}
  },
  "strengths": {
    "text1": ["strengths of text 1"],
    "text2": ["strengths of text 2"]
  },
  "weaknesses": {
    "text1": ["weaknesses of text 1"],
    "text2": ["weaknesses of text 2"]
  },
  "better_text": "text1/text2/equal",
  "recommendation": "which to use for what purpose",
  "improvement_suggestions": {
    "text1": ["how to improve text 1"],
    "text2": ["how to improve text 2"]
  }
}`
}

func buildFactualComparisonPrompt(detailLevel string) string {
	return `Compare the factual content and accuracy of these texts. Return a JSON object with:
{
  "factual_consistency": 0.0-1.0,
  "shared_facts": ["facts present in both"],
  "contradicting_facts": [
    {"fact": "description", "text1_version": "...", "text2_version": "..."}
  ],
  "unique_facts": {
    "text1": ["facts only in text 1"],
    "text2": ["facts only in text 2"]
  },
  "fact_verification_needed": ["claims requiring verification"],
  "statistical_comparison": "comparison of any statistics used",
  "source_comparison": "comparison of sources cited",
  "credibility_assessment": {
    "text1_credibility": 0.0-1.0,
    "text2_credibility": 0.0-1.0
  },
  "recommendation": "which text is more factually reliable"
}`
}

func buildComprehensiveComparisonPrompt(detailLevel string) string {
	return `Perform a comprehensive comparison of these two texts. Return a detailed JSON object with:
{
  "overall_similarity": 0.0-1.0,
  "relationship": "identical/versions/related/unrelated",
  "content_comparison": {
    "shared_topics": ["common topics"],
    "unique_to_text1": ["unique topics in text 1"],
    "unique_to_text2": ["unique topics in text 2"],
    "content_overlap": 0.0-1.0
  },
  "style_comparison": {
    "text1_style": "style description",
    "text2_style": "style description",
    "style_similarity": 0.0-1.0
  },
  "quality_comparison": {
    "text1_quality": 0.0-1.0,
    "text2_quality": 0.0-1.0,
    "better_text": "text1/text2/equal"
  },
  "key_differences": ["main differences"],
  "key_similarities": ["main similarities"],
  "statistical_summary": {
    "text1": {"words": n, "sentences": n},
    "text2": {"words": n, "sentences": n}
  },
  "use_case_recommendations": {
    "text1_best_for": "ideal use cases for text 1",
    "text2_best_for": "ideal use cases for text 2"
  },
  "verdict": "overall comparison summary"
}`
}

func NewCompareTextsAction() sdk.Action {
	return &CompareTextsAction{}
}
