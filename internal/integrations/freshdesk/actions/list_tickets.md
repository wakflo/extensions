# List Tickets

## Description

Retrieve a list of tickets based on filter criteria.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| filter | string | No | Filter tickets by predefined filters (all_tickets, open, pending, resolved, closed, new_and_my_open, watching, deleted) |
| page | number | No | Page number for pagination (default: 1) |
| per_page | number | No | Number of results per page (default: 30, max: 100) |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Sample Response

```json
{
  "data": [
    {
      "id": 123,
      "subject": "Sample Ticket 1",
      "status": 2,
      "priority": 1,
      "created_at": "2023-12-01T12:30:45Z"
    },
    {
      "id": 124,
      "subject": "Sample Ticket 2",
      "status": 3,
      "priority": 2,
      "created_at": "2023-12-02T10:15:30Z"
    }
  ]
}
```