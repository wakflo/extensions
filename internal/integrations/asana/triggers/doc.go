package triggers

import _ "embed"

//go:embed task_created.md
var taskCreatedDocs string

//go:embed task_updated.md
var taskUpdatedDocs string

//go:embed task_completed.md
var taskCompletedDocs string
