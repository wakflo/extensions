package openai

import (
	_ "embed"
)

//go:embed docs/openai.mdx
var docs string

//go:embed docs/openai.mdx
var breakLoopDocs string
