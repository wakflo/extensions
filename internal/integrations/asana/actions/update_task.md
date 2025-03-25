# Update Task
## Description
Updates an existing Asana task with new details such as name, notes, and completion status.

## Properties
| Name | Type | Required | Description |
|------|------|----------|-------------|
| Project | String | Yes | The project containing the task to update |
| Task | String | Yes | The task to update |
| Name | String | No | New name for the task |
| Notes | Text | No | New notes or description for the task |
| Completed | Boolean | No | Mark the task as completed or incomplete |

## Details
- **Type**: sdkcore.ActionTypeNormal

## Usage
This action updates an existing Asana task. You must select a project and task, then provide at least one field to update. You can modify the task's name, notes, and completion status.

The task dropdown dynamically shows only tasks from the selected project, making it easier to find the specific task you want to update.

## Example Response
```json
{
  "data": {
    "gid": "1234567890",
    "name": "Updated Task Name",
    "resource_type": "task",
    "completed": true,
    "completed_at": "2023-06-16T14:22:43.000Z",
    "created_at": "2023-06-15T10:30:00.000Z",
    "modified_at": "2023-06-16T14:22:43.000Z",
    "notes": "This task has been updated with new information.",
    "projects": [
      {
        "gid": "987654321",
        "name": "Project Name",
        "resource_type": "project"
      }
    ]
  }
}
```