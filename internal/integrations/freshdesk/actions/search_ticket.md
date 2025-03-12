# Search Tickets

## Description

Search for tickets using various parameters like keywords, statuses, and priorities.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| query | string | Yes | Search query string (e.g., 'status:open priority:high') |
| page | number | No | Page number for pagination (default: 1) |
| per_page | number | No | Number of results per page (default: 30, max: 100) |
| filter_by| string | No | Optional filter to narrow down search results |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Sample Response

```json
{
  "results": [
    {
      "id": 123,
      "subject": "Search Result Ticket 1",
      "status": 2,
      "priority": 3,
      "created_at": "2023-12-01T12:30:45Z"
    },
    {
      "id": 125,
      "subject": "Search Result Ticket 2",
      "status": 2,
      "priority": 3,
      "created_at": "2023-12-03T09:45:20Z"
    }
  ],
  "total": 2
}
```