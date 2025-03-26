# Update Task

## Description

Update an existing task in Wrike with new properties such as status.

## Properties

| Name        | Type     | Description                                                      | Required | Default |
|-------------|----------|------------------------------------------------------------------|----------|---------|
| taskId      | string   | The ID of the Wrike task to update                               | Yes      | -       |
| title       | string   | The new title of the task                                        | No       | -       |
| description | string   | The new detailed description of the task                         | No       | -       |
| status      | string   | The new status of the task (Active, Completed, Deferred, Cancelled) | No    | -       |
| importance  | string   | The new importance level of the task (High, Normal, Low)         | No       | -       |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Example Response

```json
{
  "id": "IEADTSKYA5CKABNW",
  "accountId": "IEADTSKY",
  "title": "Updated Task",
  "description": "This task has been updated via the API",
  "briefDescription": "Updated task",
  "parentIds": ["IEADTSKYA5CKAARW"],
  "superParentIds": ["IEADTSKYA5CKAARW"],
  "scope": "WsFolder",
  "status": "Completed",
  "importance": "High",
  "createdDate": "2023-03-20T14:30:45.000Z",
  "updatedDate": "2023-03-21T10:15:22.000Z",
  "completedDate": "2023-03-21T10:15:22.000Z",
  "dates": {
    "type": "Planned",
    "start": "2023-03-21T09:00:00.000Z",
    "due": "2023-03-25T18:00:00.000Z"
  },
  "responsibleIds": ["KUAIJTSKJA", "KUAIJTSKLM"],
  "authorIds": ["KUAIJTSKJA"]
}
```