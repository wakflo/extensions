package actions

import _ "embed"

//go:embed chat_gemini.md
var chatGeminiDocs string

//go:embed embeddings_gemini.md
var embeddingsGeminiDocs string

//go:embed translate_text_gemini.md
var translateTextGeminiDocs string

//go:embed summarize_text_gemini.md
var summarizeTextGeminiDocs string

//go:embed function_calling_gemini.md
var functionCallingGeminiDocs string

//go:embed analyze_image_gemini.md
var analyzeImageGeminiDocs string
