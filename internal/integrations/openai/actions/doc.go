package actions

import _ "embed"

//go:embed chat_openai.md
var chatOpenAIDocs string

//go:embed embeddings_openai.md
var embeddingsOpenAIDocs string

//go:embed data_extractor_openai.md
var dataExtractorOpenAIDocs string

//go:embed vision_openai.md
var VisionOpenAIDocs string

//go:embed image_generation_openai.md
var ImageGenerationOpenAIDocs string
