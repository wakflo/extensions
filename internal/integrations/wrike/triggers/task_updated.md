# Task Updated

## Description

Triggers when a task is updated in your Wrike account.

## Properties

| Name      | Type   | Description                                                            | Required | Default |
|-----------|--------|------------------------------------------------------------------------|----------|---------|
| folderId  | string | The ID of the Wrike folder to monitor for updated tasks                | No       | -       |
| limit     | number | Maximum number of tasks to return when triggered (1-100)               | No       | 25      |

## Details

- **Type**: sdkcore.TriggerTypePolling
- This trigger polls the Wrike API every 5 minutes to check for updated tasks.
- The trigger compares the update date of tasks against the last time the trigger was run.
- It filters out tasks that were just created and not actually updated by comparing the creation and update dates.
- When run for the first time, it will retrieve tasks updated in the last hour.

## Example Response

```json
[
  {
    "id": "IEADTSKYA5CKABNW",
    "accountId": "IEADTSKY",
    "title": "Updated Task",
    "description": "This task has been updated in Wrike",
    "status": "Completed",
    "importance": "High",
    "createdDate": "2023-03-15T09:45:30.000Z",
    "updatedDate": "2023-03-20T14:30:45.000Z",
    "completedDate": "2023-03-20T14:30:45.000Z",
    "briefDescription": "Updated task",
    "parentIds": ["IEADTSKYA5CKAARW"],
    "superParentIds": ["IEADTSKYA5CKAARW"],
    "scope": "WsFolder",
    "authorIds": ["KUAIJTSKJA"],
    "hasAttachments": false
  }
]
```