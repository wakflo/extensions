# Create Task
## Description
Creates a new task in Asana with specified details such as name, workspace, and project assignment.

## Properties
| Name | Type | Required | Description |
|------|------|----------|-------------|
| Task Name | String | Yes | The name of the task to create |
| Workspace | String | No | The Asana workspace where the task will be created |
| Project | String | No | A project to add the newly created task to |

## Details
- **Type**: sdkcore.ActionTypeNormal

## Usage
This action creates a new task in Asana. You must provide a task name, and can optionally specify a workspace and a project to add the task to. If both workspace and project are specified, the task will be created in the workspace and added to the project.

## Example Response
```json
{
  "data": {
    "gid": "1234567890",
    "name": "Complete Q3 report",
    "resource_type": "task",
    "completed": false,
    "completed_at": null,
    "created_at": "2023-06-15T10:30:00.000Z",
    "due_at": null,
    "due_on": null,
    "notes": "",
    "projects": [
      {
        "gid": "987654321",
        "name": "Content Calendar",
        "resource_type": "project"
      }
    ],
    "workspace": {
      "gid": "543216789",
      "name": "Marketing",
      "resource_type": "workspace"
    }
  }
}
```