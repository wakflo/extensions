# Create Task

## Description
Creates a new task in a specified ClickUp list with customizable details like name, description, priority, and assignees.

## Details
- **Type**: sdkcore.ActionTypeNormal

## Properties
| Name | Type | Required | Description |
|------|------|----------|-------------|
| workspace-id | String | Yes | Select a workspace. |
| space-id | String | Yes | Select a space. |
| folder-id | String | Yes | Select a folder. |
| list-id | String | Yes | Select a list to create task in. |
| assignee-id | String | Yes | Select an assignee. |
| name | String | Yes | The name of the task. |
| description | String | No | The description of the task. |
| priority | String | No | The priority level of the task. |

## Sample Output
```json
{
  "id": "abc123",
  "name": "Example Task",
  "description": "This is a sample task",
  "status": {
    "status": "Open",
    "color": "#d3d3d3"
  },
  "priority": {
    "priority": "High", 
    "color": "#f50000"
  },
  "date_created": "1647354847362",
  "date_updated": "1647354847362",
  "creator": {
    "id": "123456",
    "username": "John Doe",
    "email": "john@example.com"
  }
}
```