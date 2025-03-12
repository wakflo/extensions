# Task Created

## Description

The "Task Created" trigger monitors your HubSpot CRM for newly created tasks. When a new task is detected, this trigger will fire and allow you to initiate a series of automated workflows, such as sending notifications, creating follow-up actions, or syncing with other systems.

## Properties

| Name       | Type   | Required | Description                                                                                        |
|------------|--------|----------|----------------------------------------------------------------------------------------------------|
| properties | string | No       | Comma-separated list of task properties to include in the response (e.g., hs_task_subject,hs_task_body,hs_task_priority) |

## Details

- **Type**: sdkcore.TriggerTypePolling

## Output

This trigger outputs a list of newly created tasks with their properties. The output structure will include:

```json
{
  "results": [
    {
      "id": "456",
      "properties": {
        "hs_task_subject": "Follow up with customer",
        "hs_task_body": "Call customer to discuss recent support ticket",
        "hs_task_priority": "HIGH",
        "hs_createdate": "2024-03-10T15:45:30.456Z"
      },
      "createdAt": "2024-03-10T15:45:30.456Z",
      "updatedAt": "2024-03-10T15:45:30.456Z"
    }
  ]
}
```

## Notes

- The trigger polls the HubSpot API for tasks that have been created since the last time the trigger ran.
- If you don't specify any properties, the response will include default task properties.
- You can optimize performance by requesting only the specific properties you need.
- HubSpot API has rate limits, so be mindful of how frequently this trigger runs.
- The trigger will return up to 100 newly created tasks per run.