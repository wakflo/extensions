package actions

import _ "embed"

//go:embed create_project.md
var createProjectDocs string

//go:embed create_task.md
var createTaskDocs string

//go:embed get_active_task.md
var getActiveTaskDocs string

//go:embed get_project.md
var getProjectDocs string

//go:embed list_project_collaborators.md
var listProjectCollaboratorsDocs string

//go:embed list_projects.md
var listProjectsDocs string

//go:embed list_task.md
var listTaskDocs string

//go:embed update_task.md
var updateTaskDocs string

//go:embed update_project.md
var updateProjectDocs string
