# Task Created

## Description

Triggers when a new task is created in your Wrike account.

## Properties

| Name      | Type   | Description                                                            | Required | Default |
|-----------|--------|------------------------------------------------------------------------|----------|---------|
| folderId  | string | The ID of the Wrike folder to monitor for new tasks                    | No       | -       |
| limit     | number | Maximum number of tasks to return when triggered (1-100)               | No       | 25      |

## Details

- **Type**: sdkcore.TriggerTypePolling
- This trigger polls the Wrike API every 5 minutes to check for new tasks.
- The trigger compares the creation date of tasks against the last time the trigger was run.
- When run for the first time, it will retrieve tasks created in the last hour.

## Example Response

```json
[
  {
    "id": "IEADTSKYA5CKABNW",
    "accountId": "IEADTSKY",
    "title": "New Task",
    "description": "This is a new task created in Wrike",
    "status": "Active",
    "importance": "Normal",
    "createdDate": "2023-03-20T14:30:45.000Z",
    "updatedDate": "2023-03-20T14:30:45.000Z",
    "briefDescription": "New task",
    "parentIds": ["IEADTSKYA5CKAARW"],
    "superParentIds": ["IEADTSKYA5CKAARW"],
    "scope": "WsFolder",
    "authorIds": ["KUAIJTSKJA"],
    "hasAttachments": false
  }
]
```