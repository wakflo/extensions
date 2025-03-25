# Task Updated

## Description

Triggers a workflow when an existing task is updated in ClickUp, including changes to status, priority, assignees, or due dates, enabling automated reactions to task modification events.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| list-id | String | Conditional | The ID of the ClickUp list to monitor for updated tasks. Required if Workspace ID is not provided. |
| workspace-id | String | Conditional | The ID of the ClickUp workspace to monitor for updated tasks. Required if List ID is not provided. |
| status | String | No | Only trigger when tasks are updated to this status. Leave empty to trigger on any status change. |

## Sample Output

```json
{
  "id": "abc123",
  "name": "Updated Task",
  "status": {
    "status": "In Progress",
    "color": "#4194f6"
  },
  "date_created": "1647354847362",
  "date_updated": "1647354987362",
  "update_fields": [
    "status",
    "assignees",
    "due_date"
  ],
  "assignees": [
    {
      "id": 123456,
      "username": "John Doe",
      "email": "john@example.com"
    }
  ]
}
```

## Notes

- This trigger uses polling to check for updated tasks at regular intervals.
- Only tasks updated after the last poll will be returned.
- You must provide either a List ID or a Workspace ID.
- When using a Workspace ID, the trigger will detect tasks updated across all lists in the workspace.
- The trigger filters out newly created tasks to focus only on genuine updates to existing tasks.
- You can optionally filter by status to only trigger when tasks are updated to a specific status.