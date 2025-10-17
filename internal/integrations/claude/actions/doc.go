package actions

import _ "embed"

//go:embed chat_claude.md
var chatClaudeDocs string

//go:embed chat_with_image.md
var chatImageDocs string

//go:embed translate_text.md
var translateTextDocs string

//go:embed summarize_text.md
var summarizeTetDocs string

//go:embed compare_text.md
var compareTextsDocs string

//go:embed analyze_text.md
var analyzeTextDocs string
