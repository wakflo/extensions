# List Tasks

## Description

Retrieve a list of tasks from Wrike with optional filtering parameters.

## Properties

| Name       | Type   | Description                                                           | Required | Default      |
|------------|--------|-----------------------------------------------------------------------|----------|--------------|
| folderId   | string | The ID of the Wrike folder containing the tasks to list               | No       | -            |
| status     | string | Filter tasks by status (Active, Completed, Deferred, Cancelled, All)  | No       | Active       |
| limit      | number | Maximum number of tasks to return (1-100)                             | No       | 25           |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Example Response

```json
[
  {
    "id": "IEADTSKYA5CKABJN",
    "accountId": "IEADTSKY",
    "title": "Example Task 1",
    "description": "This is an example task",
    "status": "Active",
    "importance": "Normal",
    "createdDate": "2023-03-15T10:30:45.000Z",
    "updatedDate": "2023-03-15T15:20:15.000Z"
  },
  {
    "id": "IEADTSKYA5CKABJ2",
    "accountId": "IEADTSKY",
    "title": "Example Task 2",
    "description": "This is another example task",
    "status": "Active",
    "importance": "High",
    "createdDate": "2023-03-16T09:15:30.000Z",
    "updatedDate": "2023-03-16T14:45:22.000Z"
  }
]
```