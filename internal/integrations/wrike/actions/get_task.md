# Get Task

## Description

Retrieves detailed information about a specific task in Wrike by its ID.

## Properties

| Name    | Type   | Description                           | Required |
|---------|--------|---------------------------------------|----------|
| taskId  | string | The ID of the Wrike task to retrieve  | Yes      |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Example Response

```json
{
  "id": "IEADTSKYA5CKABJN",
  "accountId": "IEADTSKY",
  "title": "Example Task",
  "description": "This is an example task",
  "briefDescription": "Example task",
  "parentIds": ["IEADTSKYA5CKAARW"],
  "superParentIds": ["IEADTSKYA5CKAARW"],
  "scope": "WsFolder",
  "status": "Active",
  "importance": "Normal",
  "createdDate": "2023-03-15T10:30:45.000Z",
  "updatedDate": "2023-03-15T15:20:15.000Z",
  "dates": {
    "type": "Planned",
    "duration": 86400,
    "start": "2023-03-16T09:00:00.000Z",
    "due": "2023-03-17T18:00:00.000Z"
  },
  "responsibleIds": ["KUAIJTSKJA"],
  "authorIds": ["KUAIJTSKJA"],
  "hasAttachments": false,
  "priority": "Normal"
}
```