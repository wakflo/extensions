# Update Ticket

## Description

Update the properties and fields of an existing Freshdesk ticket.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| ticket_id | number | Yes | The ID of the ticket to update |
| subject | string | No | The updated subject of the ticket |
| description | string | No | The updated description of the ticket |
| priority | number | No | The updated priority level (1-4, where 1=Low, 2=Medium, 3=High, 4=Urgent) |
| status | number | No | The updated status (2=Open, 3=Pending, 4=Resolved, 5=Closed) |
| type | string | No | The updated type of the ticket |
| tags | string | No | Comma-separated list of updated tags |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Sample Response

```json
{
  "id": 123,
  "subject": "Updated Ticket Subject",
  "description": "This ticket has been updated via API",
  "status": 3,
  "priority": 2,
  "requester_id": 456,
  "updated_at": "2023-12-01T15:45:30Z"
}
```