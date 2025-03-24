# Get Task
## Description
Retrieves detailed information for a specific Asana task.

## Properties
| Name | Type | Required | Description |
|------|------|----------|-------------|
| Project | String | Yes | The project containing the task |
| Task | String | Yes | The task to retrieve details for |

## Details
- **Type**: sdkcore.ActionTypeNormal

## Usage
This action retrieves comprehensive details for a specific task in Asana. You must first select a project, which will then populate the task dropdown with tasks from that project.

The task dropdown dynamically updates based on the selected project, making it easier to find the specific task you're looking for.

## Example Response
```json
{
  "data": {
    "gid": "1234567890",
    "name": "Finalize marketing campaign",
    "resource_type": "task",
    "completed": false,
    "completed_at": null,
    "created_at": "2023-06-12T09:15:00.000Z",
    "modified_at": "2023-06-14T16:30:22.000Z",
    "due_at": null,
    "due_on": "2023-06-30",
    "notes": "Campaign should target our primary demographic with focus on new product features.",
    "assignee": {
      "gid": "987654321",
      "name": "Alex Johnson",
      "resource_type": "user"
    },
    "assignee_status": "upcoming",
    "projects": [
      {
        "gid": "5678901234",
        "name": "Q2 Marketing",
        "resource_type": "project"
      }
    ],
    "parent": null,
    "workspace": {
      "gid": "2468013579",
      "name": "Marketing Team",
      "resource_type": "workspace"
    }
  }
}
```