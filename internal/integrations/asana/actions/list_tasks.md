# List Tasks
## Description
Retrieves a list of tasks from a specific Asana project.

## Properties
| Name | Type | Required | Description |
|------|------|----------|-------------|
| Project | String | Yes | The project to list tasks from |
| Limit | Number | No | Maximum number of tasks to return (default: 50) |

## Details
- **Type**: sdkcore.ActionTypeNormal

## Usage
This action retrieves all tasks from a specified Asana project. You can optionally limit the number of results returned.

The action returns comprehensive task details including names, assignees, due dates, and completion status for each task in the project.

## Example Response
```json
{
  "data": [
    {
      "gid": "1234567890",
      "name": "Complete project proposal",
      "resource_type": "task",
      "completed": false,
      "due_on": "2023-06-20",
      "notes": "Include all requirements and timeline",
      "assignee": {
        "gid": "987654321",
        "name": "Jane Smith",
        "resource_type": "user"
      }
    },
    {
      "gid": "2345678901",
      "name": "Schedule kickoff meeting",
      "resource_type": "task",
      "completed": true,
      "due_on": "2023-06-10",
      "notes": "",
      "assignee": null
    },
    {
      "gid": "3456789012",
      "name": "Create initial wireframes",
      "resource_type": "task",
      "completed": false,
      "due_on": "2023-06-25",
      "notes": "Focus on main user flows",
      "assignee": {
        "gid": "987654321",
        "name": "Jane Smith",
        "resource_type": "user"
      }
    }
  ]
}
```