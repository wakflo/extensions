package actions

import _ "embed"

//go:embed get_task.md
var getTaskDocs string

//go:embed list_tasks.md
var listTasksDocs string

//go:embed create_task.md
var createTaskDocs string

//go:embed update_task.md
var updateTaskDocs string
