# Create Task

## Description

Create a new task in Wrike with specified title, description, and status.

## Properties

| Name        | Type     | Description                                                      | Required | Default |
|-------------|----------|------------------------------------------------------------------|----------|---------|
| folderId    | string   | The ID of the Wrike folder where the task will be created        | Yes      | -       |
| title       | string   | The title of the task                                            | Yes      | -       |
| description | string   | The detailed description of the task                             | No       | -       |
| status      | string   | The status of the task (Active, Completed, Deferred, Cancelled)  | No       | Active  |
| importance  | string   | The importance level of the task (High, Normal, Low)             | No       | Normal  |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Example Response

```json
{
  "id": "IEADTSKYA5CKABNW",
  "accountId": "IEADTSKY",
  "title": "New Task",
  "description": "This is a new task created via the API",
  "briefDescription": "New task",
  "parentIds": ["IEADTSKYA5CKAARW"],
  "superParentIds": ["IEADTSKYA5CKAARW"],
  "scope": "WsFolder",
  "status": "Active",
  "importance": "Normal",
  "createdDate": "2023-03-20T14:30:45.000Z",
  "updatedDate": "2023-03-20T14:30:45.000Z",
  "dates": {
    "type": "Planned",
    "start": "2023-03-21T09:00:00.000Z",
    "due": "2023-03-25T18:00:00.000Z"
  },
  "responsibleIds": ["KUAIJTSKJA"],
  "authorIds": ["KUAIJTSKJA"]
}
```