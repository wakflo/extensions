# Get Ticket

## Description

Retrieve detailed information about a specific ticket by its ID.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| ticket_id | number | Yes | The ID of the ticket to retrieve |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Sample Response

```json
{
  "id": 123,
  "subject": "Sample Ticket",
  "description": "This is a sample ticket details",
  "status": 2,
  "priority": 1,
  "requester_id": 456,
  "responder_id": 789,
  "created_at": "2023-12-01T12:30:45Z", 
  "updated_at": "2023-12-01T14:20:15Z",
  "due_by": "2023-12-03T12:30:45Z",
  "fr_due_by": "2023-12-02T12:30:45Z",
  "tags": ["support", "urgent"]
}
```
