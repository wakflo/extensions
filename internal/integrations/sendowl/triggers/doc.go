// triggers/doc.go
package triggers

import _ "embed"

//go:embed order_completed.md
var orderCompletedDocs string

//go:embed product_updated.md
var productUpdatedDocs string
