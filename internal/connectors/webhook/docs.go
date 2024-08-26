package webhook

import (
	_ "embed"
)

//go:embed docs/webhook.mdx
var webhookDocs string

//go:embed docs/trigger_catch.mdx
var triggerCatchDocs string
