# Task Created

## Description

Triggers a workflow when a new task is created in a specified ClickUp workspace or list, allowing you to automate subsequent actions based on new task creation events.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| list-id | String | Conditional | The ID of the ClickUp list to monitor for new tasks. Required if Workspace ID is not provided. |
| workspace-id | String | Conditional | The ID of the ClickUp workspace to monitor for new tasks. Required if List ID is not provided. |

## Sample Output

```json
{
  "id": "abc123",
  "name": "New Task",
  "status": {
    "status": "Open",
    "color": "#d3d3d3"
  },
  "date_created": "1647354847362",
  "creator": {
    "id": 123456,
    "username": "John Doe",
    "email": "john@example.com"
  }
}
```

## Notes

- This trigger uses polling to check for new tasks at regular intervals.
- Only tasks created after the last poll will be returned.
- You must provide either a List ID or a Workspace ID.
- When using a Workspace ID, the trigger will detect tasks created across all lists in the workspace.