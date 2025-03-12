package triggers

import _ "embed"

//go:embed contact_updated.md
var contactUpdatedDoc string

//go:embed ticket_updated.md
var ticketCreatedDoc string

//go:embed deal_updated.md
var dealUpdatedDoc string

//go:embed task_created.md
var taskCreatedDoc string
